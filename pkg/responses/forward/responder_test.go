package forward

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestNewResponderHasTheExpectedValues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/path", r.URL.Path)
		rw.WriteHeader(202)
		rw.Write([]byte("{\"success\": true}"))
	}))
	defer srv.Close()

	req, _ := http.NewRequest("PUT", "http://example.com/path", nil)
	serverURL, _ := url.Parse(srv.URL)
	rspnr := NewResponder(serverURL.Host)
	defer rspnr.Close()

	res, err := rspnr.Respond(req)
	assert.Nil(t, err)
	assert.Equal(t, 202, res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body")
	}

	assert.Equal(t, "{\"success\": true}", string(resBody))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
