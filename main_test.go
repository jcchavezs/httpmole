package main

import "testing"

func TestUnescapeBodyRemovesWrappingQuotes(t *testing.T) {
	body := []byte(`"{}"`)
	if want, have := `{}`, string(unescapeBody(body)); want != have {
		t.Errorf("failed to scape body, want %q, have %q", want, have)
	}
}
