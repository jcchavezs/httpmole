package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jcchavezs/httpmole/pkg/flags"
	"github.com/jcchavezs/httpmole/pkg/format"
	"github.com/jcchavezs/httpmole/pkg/responses"
	"github.com/jcchavezs/httpmole/pkg/responses/file"
	"github.com/jcchavezs/httpmole/pkg/responses/forward"
	"github.com/jcchavezs/httpmole/pkg/responses/static"
)

func main() {
	var (
		port           int
		resFilepath    string
		resStatusCode  int
		resHeaderLines flags.Slice
		resFrom        string
		showResponse   bool
		durationInMS   int
	)

	flag.IntVar(&port, "p", 10080, "Listening port")
	flag.StringVar(
		&resFilepath,
		"response-file",
		"",
		"Response filepath in JSON format. See https://github.com/jcchavezs/httpmole/blob/master/examples/response-file.json",
	)
	flag.IntVar(&resStatusCode, "response-status", 200, "Response status code")
	flag.Var(&resHeaderLines, "response-header", "Response headers e.g. location:/login")
	flag.StringVar(
		&resFrom,
		"response-from",
		"",
		"Response source hostport, e.g. realservice:1234",
	)
	flag.BoolVar(&showResponse, "show-response", false, "Display the response along with the request")
	flag.IntVar(&durationInMS, "duration-ms", 0, "Duration of the operation in milliseconds")
	flag.Parse()

	var resp responses.Responder
	if resFrom != "" {
		resp = forward.NewResponder(resFrom)
	} else if resFilepath != "" {
		resp = file.NewResponder(resFilepath)
	} else {
		resp = static.NewResponder(resStatusCode, *toHeadersMap(resHeaderLines))
	}
	defer resp.Close()

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		logRequest(req, os.Stdout)

		var currentDurationInMS = durationInMS

		if currentDurationInMS > 0 {
			time.Sleep(time.Duration(durationInMS) * time.Millisecond)
		}

		var (
			res *http.Response
			err error
		)

		if strings.HasPrefix(req.URL.Path, "/proxy/") {
			hostport, newReq, ok := newProxyRequest(req)
			if ok {
				currentTime := time.Now().UnixMilli()
				res, err = forward.NewResponder(hostport).Respond(newReq)
				currentDurationInMS = int(time.Now().UnixMilli() - currentTime)
			} else {
				res, err = resp.Respond(req)
			}
		} else {
			res, err = resp.Respond(req)
		}

		if err == nil {
			var logWriter io.Writer
			if showResponse {
				logWriter = os.Stdout
			}
			if currentDurationInMS > 0 {
				res.Header.Set("Server-Timing", fmt.Sprintf("app;dur=%.2f", float64(currentDurationInMS/1000.0)))
			}
			writeResponse(res, rw, logWriter)
		} else {
			log.Printf("failed to resolve the response: %v\n\n", err)
			rw.WriteHeader(502)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func logRequest(r *http.Request, w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%s %s %s", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL.String())))
	for k, v := range r.Header {
		w.Write([]byte(fmt.Sprintf("\n > %s: %v", k, strings.Join(v, "; "))))
	}

	if r.Method != "GET" {
		rBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read request body: %v", err)
			return
		}
		r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(rBody))
		if len(rBody) > 0 {
			w.Write([]byte(fmt.Sprintf("\n\n%s", string(rBody))))
		}
	}

	w.Write([]byte("\n\n"))
}

func writeResponse(res *http.Response, w http.ResponseWriter, lw io.Writer) {
	if lw != nil {
		lw.Write([]byte(fmt.Sprintf("Status Code: %d", res.StatusCode)))
	}

	for k, v := range res.Header {
		w.Header().Add(k, strings.Join(v, "; "))
		if lw != nil {
			lw.Write([]byte(fmt.Sprintf("\n > %s: %v", k, strings.Join(v, "; "))))
		}
	}
	w.WriteHeader(res.StatusCode)

	defer res.Body.Close()

	var (
		body []byte
		err  error
	)

	if res.StatusCode != http.StatusNoContent {
		body, err = io.ReadAll(res.Body)
		if err == nil && len(body) > 0 {
			// Here we assume that whatever response adapter is adding a best guess content type
			// and use that information for further tweaks.
			formatter := format.GetFormatterContentType(res.Header.Get("Content-Type"))
			w.Write(formatter(body, format.Expanded))
			if lw != nil {
				lw.Write([]byte(fmt.Sprintf("\n\n%s", string(formatter(body, format.Minified)))))
			}
		}
	}

	if lw != nil {
		lw.Write([]byte("\n\n"))
	}
}

func toHeadersMap(headersLine []string) *http.Header {
	headers := &http.Header{}

	if len(headersLine) != 0 {
		for _, headerLine := range headersLine {
			headerLinePieces := strings.SplitN(headerLine, ":", 2)
			headers.Add(headerLinePieces[0], headerLinePieces[1])
		}
	}
	return headers
}

func newProxyRequest(req *http.Request) (string, *http.Request, bool) {
	if !strings.HasPrefix(req.URL.Path, "/proxy/") {
		panic("newProxyRequest should not be called with a request that already has a /proxy/ prefix")
	}

	hostport, newPath, _ := strings.Cut(req.URL.Path[7:], "/")
	if hostport == "" {
		return "", nil, false
	}

	newReq := req.Clone(req.Context())
	newReq.URL.Path = "/" + newPath
	return hostport, newReq, true
}
