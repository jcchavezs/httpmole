package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToHeadersMap(t *testing.T) {
	h := toHeadersMap([]string{"x-request-id:abc123", "location:http://localhost:8080/login"})
	require.Equal(t, "abc123", h.Get("x-request-id"))
	require.Equal(t, "http://localhost:8080/login", h.Get("location"))
}
