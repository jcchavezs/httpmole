package file

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	fpath "path/filepath"

	"github.com/jcchavezs/httpmole/pkg/responses"

	"gopkg.in/fsnotify.v1"
)

type responder struct {
	mustSyncResponse bool
	watcher          *fsnotify.Watcher
	mutex            sync.Mutex
	response         *response
	filepath         string
}

// NewResponder returns a Responder that responds with a dynamic response
// being specified in a given filepath.
func NewResponder(filepath string) responses.Responder {
	var err error
	filepath, err = fpath.Abs(filepath)
	if err != nil {
		log.Fatalf("failed to stablish absolute path for %q: %v", filepath, err)
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatalf("failed to check response file: %v", err)
	}
	r := &responder{filepath: filepath, mustSyncResponse: true, response: &response{}}
	r.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("failed to start the response watcher: %v", err)
	}

	err = r.watcher.Add(filepath)
	if err == nil {
		go r.checkNewResponse()
	} else {
		log.Fatalf("failed to add the file watcher for %q: %v", filepath, err)
	}

	return r
}

func (fr *responder) Respond(_ *http.Request) (*http.Response, error) {
	if fr.mustSyncResponse {
		if err := fr.loadResponse(); err != nil {
			log.Fatalf("failed to load response: %v", err)
		}
	}
	res := &http.Response{
		StatusCode: fr.response.statusCode,
		Header:     fr.response.headers,
		Body:       ioutil.NopCloser(nil),
	}

	if len(fr.response.body) > 0 {
		res.Body = ioutil.NopCloser(bytes.NewBuffer(fr.response.body))
	}

	return res, nil
}

func (fr *responder) Close() {
	fr.watcher.Close()
}

func (fr *responder) checkNewResponse() {
	for {
		select {
		case <-fr.watcher.Events:
			fr.mutex.Lock()
			fr.mustSyncResponse = true
			fr.mutex.Unlock()
		case err := <-fr.watcher.Errors:
			if err == nil {
				return
			}
			log.Fatalf("failed to watch response file: %v", err)
		}
	}
}

// loadResponse reads the resFilepath and overrides the values in the provided
// response. If the json parsing fails, the response won't be overwriten.
func (fr *responder) loadResponse() error {
	data, err := ioutil.ReadFile(fr.filepath)
	if err != nil {
		return err
	}
	readingRes := response{}
	err = json.Unmarshal(data, &readingRes)
	if err != nil {
		return err
	}

	fr.response.copyFrom(readingRes)
	return nil
}
