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
		isFileRead, err := fr.loadResponse()
		if !isFileRead {
			log.Fatalf("failed to load response: %v", err)
		} else if err != nil {
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
	return (body[0] == '{' || body[0] == '[')
}

func unescapeQuotesInBody(body []byte) []byte {
	if len(body) == 0 {
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
func (fr *responder) loadResponse() (bool, error) {
	data, err := ioutil.ReadFile(fr.filepath)
	if err != nil {
		return false, err
	}
	readingRes := response{}
	err = json.Unmarshal(data, &readingRes)
	if err != nil {
		return true, err
	}

	if err = readingRes.validate(); err != nil {
		return true, err
	}

	if isRawJSON(readingRes.body) {
		if readingRes.headers.Get("Content-Type") == "" {
			readingRes.headers.Set("Content-Type", "application/json")
		}
	} else {
		readingRes.body = unescapeQuotesInBody(readingRes.body)
	}

	fr.response.copyFrom(readingRes)
	return true, nil
}
