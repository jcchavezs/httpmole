package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	r := &responder{filepath: filepath, mustSyncResponse: true, response: &response{}}

	return r
}

func (fr *responder) lazyLoadFileWatcher() error {
	var err error
	fr.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("failed to start the response watcher: %v", err)
	}

	err = fr.watcher.Add(fr.filepath)
	if err == nil {
		go fr.checkNewResponse()
	} else {
		return fmt.Errorf("failed to add the file watcher for %q: %v", fr.filepath, err)
	}

	return nil
}

func (fr *responder) Respond(_ *http.Request) (*http.Response, error) {
	if fr.watcher == nil {
		// the service should be still be up even if the fail does not exist, mostly because
		// the usual flow is to open the app and then create the file, hence the deferring of
		// the file watcher loading until the request comes.
		if err := fr.lazyLoadFileWatcher(); err != nil {
			return nil, err
		}
	}

	if fr.mustSyncResponse {
		if err := fr.loadResponse(); err != nil {
			return nil, err
		}
	}

	res := &http.Response{
		StatusCode: fr.response.statusCode,
		Header:     fr.response.headers,
		Body:       ioutil.NopCloser(bytes.NewBuffer(fr.response.body)),
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

func isRawJSON(body []byte) bool {
	if len(body) == 0 {
		return false
	}

	return (body[0] == '{' || body[0] == '[')
}

func unescapeQuotesInBody(body []byte) []byte {
	if len(body) < 2 {
		return body
	}

	// It is a string hence we need to unescape quotes
	body = body[1 : len(body)-1]

	var dstRune []rune
	strRune := []rune(string(body))
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] && strRune[i+1] == []rune{'"'}[0] {
			continue
		}
		dstRune = append(dstRune, strRune[i])
	}
	return []byte(string(dstRune))
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

	if err = readingRes.validate(); err != nil {
		return err
	}

	if isRawJSON(readingRes.body) {
		if readingRes.headers.Get("Content-Type") == "" {
			readingRes.headers.Set("Content-Type", "application/json")
		}
	} else {
		readingRes.body = unescapeQuotesInBody(readingRes.body)
	}

	fr.response.copyFrom(readingRes)
	return nil
}
