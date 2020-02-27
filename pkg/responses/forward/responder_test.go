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
		assert.Equal(t, "a0b0c123d0e0f456", r.Header.Get("x-b3-traceid"))
		rw.Header().Add("Content-Type", "application/vnd.schemaregistry.v1+json")
		rw.WriteHeader(202)
		rw.Write([]byte("{\"success\": true}"))
	}))
	defer srv.Close()

	req, _ := http.NewRequest("PUT", "http://example.com/path", nil)
	req.Header.Set("x-b3-traceid", "a0b0c123d0e0f456")
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

func TestTransformedURLHasTheExpectedValue(t *testing.T) {
	reqs := make([]*http.Request, 2)
	reqs[0] = &http.Request{
		Method:     "POST",
		URL:        &url.URL{Path: "/my/path"},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Host:       "oldhost:9411",
		RemoteAddr: "172.20.0.6:43228",
		RequestURI: "/my/path",
	}
	reqs[1], _ = http.NewRequest("PUT", "http://old_domain/my/path", nil)

	for _, req := range reqs {
		assert.Equal(t, "http://newhost:1111/my/path", transformURL(req, "newhost:1111"))
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
