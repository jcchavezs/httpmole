package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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
		port             int
		resFilepath      string
		resStatusCode    int
		resHeaderLines   flags.Slice
		resFrom          string
		logMethodMatcher *regexp.Regexp
		logPathMatcher   *regexp.Regexp
	)

	flag.IntVar(&port, "p", 10080, "Listening port")
	flag.StringVar(
		&resFilepath,
		"response-file",
		"",
		"Response filepath in JSON format. See https://github.com/jcchavezs/httpmole/blob/master/examples/response-file.json",
	)
	flag.IntVar(&resStatusCode, "response-status", 200, "Response status code")
	flag.Var(&resHeaderLines, "response-header", "Response headers")
	flag.StringVar(
		&resFrom,
		"response-from",
		"",
		"Response source hostport, e.g. realservice:1234",
	)
	logResponse := flag.Bool("log-response", false, "Logs the response along with the request")
	logMethodRegex := flag.String("log-filter-method", "", "Log only matching method e.g. `\"GET|DELETE\"`")
	logPathRegex := flag.String("log-filter-path", "", "Log only matching path e.g. \"/^(health)\"")
	flag.Parse()

	var err error
	if *logMethodRegex != "" {
		logMethodMatcher, err = regexp.Compile(*logMethodRegex)
		if err != nil {
			log.Fatal("failed to compile log-filter-method")
		}
	}
	if *logPathRegex != "" {
		logPathMatcher, err = regexp.Compile(*logPathRegex)
		if err != nil {
			log.Fatal("failed to compile log-filter-path")
		}
	}

	logRequestMatcher := makeLogRequestMatcher(logMethodMatcher, logPathMatcher)

	var resp responses.Responder
	if resFrom != "" {
		resp = forward.NewResponder(resFrom)
	} else if resFilepath != "" {
		resp = file.NewResponder(resFilepath)
	} else {
		resp = static.NewResponder(resStatusCode, *toHeadersMap(resHeaderLines))
	}
	defer resp.Close()

	logWriter := os.Stdout

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		shouldLogRequest := logRequestMatcher(req)
		if shouldLogRequest {
			logRequest(req, logWriter)
		}
		res, err := resp.Respond(req)
		if err == nil {
			var logResponseWriter io.Writer
			if *logResponse && shouldLogRequest {
				logResponseWriter = logWriter
			}
			writeResponse(res, rw, logResponseWriter)
		} else {
			log.Printf("failed to resolve the response: %v\n\n", err)
			rw.WriteHeader(502)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func makeLogRequestMatcher(logMethodMatcher, logPathMatcher *regexp.Regexp) func(r *http.Request) bool {
	if logMethodMatcher == nil && logPathMatcher == nil {
		return func(r *http.Request) bool {
			return true
		}
	}

	if logMethodMatcher == nil {
		return func(r *http.Request) bool {
			return logPathMatcher.MatchString(r.URL.Path)
		}
	} else if logPathMatcher == nil {
		return func(r *http.Request) bool {
			return logMethodMatcher.MatchString(r.Method)
		}
	} else {
		return func(r *http.Request) bool {
			return logPathMatcher.MatchString(r.URL.Path) && logMethodMatcher.MatchString(r.Method)
		}
	}
}

func logRequest(r *http.Request, w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%s %s %s", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL.String())))
	for k, v := range r.Header {
		w.Write([]byte(fmt.Sprintf("\n > %s: %v", k, strings.Join(v, "; "))))
	}

	if r.Method != "GET" {
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read request body: %v", err)
			return
		}
		r.Body.Close()

		r.Body = ioutil.NopCloser(bytes.NewBuffer(rBody))
		if rBody != nil && len(rBody) > 0 {
			w.Write([]byte(fmt.Sprintf("\n\n%s", string(rBody))))
		}
	}

	w.Write([]byte("\n\n"))
}

func writeResponse(res *http.Response, w http.ResponseWriter, lw io.Writer) {
	if lw != nil {
		lw.Write([]byte(fmt.Sprintf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))))
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
		body, err = ioutil.ReadAll(res.Body)
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
	headers := new(http.Header)

	if len(headersLine) != 0 {
		for _, headerLine := range headersLine {
			headerLinePieces := strings.Split(headerLine, ":")
			headers.Add(headerLinePieces[0], headerLinePieces[1])
		}
	}
	return headers
}
