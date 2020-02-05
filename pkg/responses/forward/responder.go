package forward

import (
	"fmt"
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
	fmt.Println(strings.Replace(req.URL.String(), req.URL.Host, fr.hostport, 1))
	fReq, _ := http.NewRequest(
		req.Method,
		strings.Replace(req.URL.String(), req.URL.Host, fr.hostport, 1),
		req.Body,
	)
	return fr.httpClient.Do(fReq)
}

func (fr *responder) Close() {}
