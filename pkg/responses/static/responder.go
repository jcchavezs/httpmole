package static

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jcchavezs/httpmole/pkg/responses"
)

type responder struct {
	statusCode int
	headers    http.Header
}

// NewResponder returns a responder that serves a static response defined
// on creation
func NewResponder(statusCode int, headers http.Header) responses.Responder {
	return &responder{statusCode: statusCode, headers: headers}
}

func (sr *responder) Respond(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: sr.statusCode,
		Header:     sr.headers,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func (sr *responder) Close() {}
