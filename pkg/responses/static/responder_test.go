package static

import (
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
	assert.Equal(t, 300, res.StatusCode)
	assert.Equal(t, "value", res.Header.Get("key"))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
