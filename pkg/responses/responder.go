package responses

import "net/http"

// Responder returns a response
type Responder interface {
	Respond(r *http.Request) (*http.Response, error)
	Close()
}
