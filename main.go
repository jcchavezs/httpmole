package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type response struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

// UnmarshalJSON unmarshals a JSON response file content
func (r *response) UnmarshalJSON(data []byte) error {
	ur := &struct {
		StatusCode int               `json:"status_code"`
		Headers    map[string]string `json:"headers"`
		Body       json.RawMessage   `json:"body"`
	}{}
	if err := json.Unmarshal(data, ur); err != nil {
		return err
	}
	r.statusCode = ur.StatusCode
	r.body = unescapeBody(ur.Body)
	r.headers = ur.Headers
	return nil
}

// unescapeBody removes trailing double quotes and escaping backslashes
func unescapeBody(str []byte) []byte {
	if str[0] == '"' {
		str = str[1:]
	}

	if str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}

	var dstRune []rune
	strRune := []rune(string(str))
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] && strRune[i+1] == []rune{'"'}[0] {
			continue
		}
		dstRune = append(dstRune, strRune[i])
	}
	return []byte(string(dstRune))
}

func (r *response) copyFrom(or response) {
	r.statusCode = or.statusCode
	r.body = or.body[:]
	r.headers = or.headers
}

func main() {
	var (
		port          int
		resFilepath   string
		resStatusCode int
	)

	flag.IntVar(&port, "p", 8081, "Listening port")
	flag.StringVar(
		&resFilepath,
		"response-file",
		"",
		"Response filepath in JSON format. See https://github.com/jcchavezs/httpmole/blob/master/examples/response-file.json",
	)
	flag.IntVar(&resStatusCode, "response-status", 200, "Response status code")
	flag.Parse()

	resNeedsSync := false
	res := response{statusCode: resStatusCode}
	if resFilepath != "" {
		resNeedsSync = true
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("failed to start the response watcher: %v", err)
		}
		defer watcher.Close()

		err = watcher.Add(resFilepath)
		if err == nil {
			go checkNewResponse(watcher, &resNeedsSync)
		} else {
			log.Fatalf("failed to add the file watcher for %q: %v", resFilepath, err)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if resNeedsSync {
			loadResponse(resFilepath, &res)
		}

		record := fmt.Sprintf("\n%s %s", r.Method, r.URL.String())
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read request body: %v", err)
			return
		}

		for k, v := range r.Header {
			record += fmt.Sprintf("\n > %s: %v", k, strings.Join(v, "; "))
		}

		if r.Method != "GET" && len(rBody) > 0 {
			record += fmt.Sprintf("\n%s\n", string(rBody))
		}

		for k, v := range res.headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(res.statusCode)
		if len(res.body) > 0 {
			w.Write(res.body)
		}
		log.Printf(record)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func checkNewResponse(watcher *fsnotify.Watcher, resNeedsSync *bool) {
	for {
		select {
		case <-watcher.Events:
			*resNeedsSync = true
		case err := <-watcher.Errors:
			log.Fatalf("failed to watch response file: %v", err)
		}
	}
}

// loadResponse reads the resFilepath and overrides the values in the provided
// response. If the json parsing fails, the response won't be overwriten.
func loadResponse(resFilepath string, res *response) error {
	data, err := ioutil.ReadFile(resFilepath)
	if err != nil {
		return err
	}
	readingRes := response{}
	err = json.Unmarshal(data, &readingRes)
	if err != nil {
		return err
	}

	res.copyFrom(readingRes)
	return nil
}
