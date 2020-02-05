package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestUnescapeBodyReturnsEmptySlice(t *testing.T) {
	body := []byte{}
	if want, have := "", string(unescapeJSONBody(body)); want != have {
		t.Errorf("failed to scape body, want %q, have %q", want, have)
	}
}

func TestUnescapeBodyRemovesWrappingQuotes(t *testing.T) {
	body := []byte(`"{}"`)
	if want, have := `{}`, string(unescapeJSONBody(body)); want != have {
		t.Errorf("failed to scape body, want %q, have %q", want, have)
	}
}

const responseContent1 = `
	{
		"status_code": 201,
		"headers": {
			"x-request-id": "abc123"
		},
		"body": "{\"success\": true}"
	}
`

const responseContent2 = `
	{
		"status_code": 403,
		"headers": {
			"x-request-id": "xyz987"
		},
		"body": "{\"success\": false}"
	}
`

func TestNewResponderHasTheExpectedValues(t *testing.T) {
	filepath := os.TempDir() + "/response.json"
	err := ioutil.WriteFile(filepath, []byte(responseContent1), 0755)
	if err != nil {
		t.Fatalf("failed to create config file: %v\n", err)
	}
	defer os.Remove(filepath)

	rspnr := NewResponder(filepath)
	defer rspnr.Close()

	res, err := rspnr.Respond(nil)
	assert.Nil(t, err)
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "abc123", res.Header.Get("x-request-id"))
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body")
	}
	assert.Equal(t, "{\"success\": true}", string(resBody))
	res.Body.Close()

	err = ioutil.WriteFile(filepath, []byte(responseContent2), 0755)
	if err != nil {
		t.Fatalf("failed to create config file: %v\n", err)
	}

	res, err = rspnr.Respond(nil)
	assert.Nil(t, err)
	assert.Equal(t, 403, res.StatusCode)
	assert.Equal(t, "xyz987", res.Header.Get("x-request-id"))
	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body")
	}
	assert.Equal(t, "{\"success\": false}", string(resBody))
	res.Body.Close()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
