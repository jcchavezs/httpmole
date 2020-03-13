package file

import (
	"encoding/json"
	"errors"
	"net/http"
)

type response struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func (r *response) copyFrom(or response) {
	r.statusCode = or.statusCode
	r.body = or.body[:]
	r.headers = or.headers
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
	r.headers = toMultiValueHeaders(ur.Headers)
	r.body = ur.Body
	return nil
}

func (r *response) validate() error {
	if r.statusCode < 100 || r.statusCode >= 599 {
		return errors.New("invalid status code")
	}

	return nil
}

func toMultiValueHeaders(singleValueHeaders map[string]string) http.Header {
	multiValueHeaders := http.Header{}
	for key, value := range singleValueHeaders {
		multiValueHeaders.Add(key, value)
	}
	return multiValueHeaders
}
