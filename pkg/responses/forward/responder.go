package forward

import (
	"net/http"
	"strings"

	"github.com/jcchavezs/httpmole/pkg/responses"
)

type responder struct {
	hostport   string
	httpClient *http.Client
}

// NewResponder returns a Responder that forwards the request to a given
// hostport and corresponding the response to the client
func NewResponder(hostport string) responses.Responder {
	return &responder{hostport: hostport, httpClient: http.DefaultClient}
}

func (fr *responder) Respond(req *http.Request) (*http.Response, error) {
	fReq, _ := http.NewRequest(
		req.Method,
		transformURL(req, fr.hostport),
		req.Body,
	)
	fReq.Header = req.Header
	return fr.httpClient.Do(fReq)
}

func transformURL(req *http.Request, replacerHost string) (url string) {
	if req.URL.Host == "" {
		url = "http://" + replacerHost + req.URL.String()
	} else {
		url = strings.Replace(req.URL.String(), req.URL.Host, replacerHost, 1)
	}
	return
}

func (fr *responder) Close() {}
