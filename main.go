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
	r.body = []byte(ur.Body)
	r.headers = ur.Headers
	return nil
}

func (r *response) copyFrom(or response) {
	r.statusCode = or.statusCode
	r.body = or.body
	r.headers = or.headers
}

func main() {
	var (
		port          int
		resFilepath   string
		resStatusCode int
	)

	flag.IntVar(&port, "p", 8081, "Application port")
	flag.StringVar(&resFilepath, "response-file", "", "Response filepath")
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
			go checkNewResponse(watcher, resFilepath, &resNeedsSync)
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
			record += fmt.Sprintf("\n%s", string(rBody))
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

func checkNewResponse(watcher *fsnotify.Watcher, resFilepath string, resNeedsSync *bool) {
	for {
		select {
		case <-watcher.Events:
			*resNeedsSync = true
		case err := <-watcher.Errors:
			log.Fatalf("failed to watch response file: %v", err)
		}
	}
}

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
