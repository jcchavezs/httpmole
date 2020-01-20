package main

import "testing"

func TestUnescapeBodyReturnsEmptySlice(t *testing.T) {
	body := []byte{}
	if want, have := "", string(unescapeBody(body)); want != have {
		t.Errorf("failed to scape body, want %q, have %q", want, have)
	}
}

func TestUnescapeBodyRemovesWrappingQuotes(t *testing.T) {
	body := []byte(`"{}"`)
	if want, have := `{}`, string(unescapeBody(body)); want != have {
		t.Errorf("failed to scape body, want %q, have %q", want, have)
	}
}
