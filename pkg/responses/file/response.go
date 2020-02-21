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
	r.body = unescapeJSONBody(ur.Body)
	r.headers = toMultiValueHeaders(ur.Headers)
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

func unescapeJSONBody(body []byte) []byte {
	if len(body) < 2 {
		return body
	}

	if body[0] == '"' {
		body = body[1:]
	}

	if body[len(body)-1] == '"' {
		body = body[:len(body)-1]
	}

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
