package static

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestNewResponderHasTheExpectedValues(t *testing.T) {
	headers := http.Header{}
	headers.Add("key", "value")
	rspnr := NewResponder(300, headers)
	defer rspnr.Close()

	res, err := rspnr.Respond(nil)
	assert.Nil(t, err)
	defer res.Body.Close()
	assert.Equal(t, 300, res.StatusCode)
	assert.Equal(t, "value", res.Header.Get("key"))
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Empty(t, body)
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
