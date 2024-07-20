package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToHeadersMap(t *testing.T) {
	h := toHeadersMap([]string{"x-request-id:abc123", "location:http://localhost:8080/login"})
	require.Equal(t, "abc123", h.Get("x-request-id"))
	require.Equal(t, "http://localhost:8080/login", h.Get("location"))
}

func TestNewProxyRequest(t *testing.T) {
	t.Run("proxy with no destination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/proxy/", nil)
		hostport, _, ok := newProxyRequest(req)
		require.Equal(t, "", hostport)
		require.False(t, ok)
	})

	t.Run("proxy with no destination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/proxy", nil)
		hostport, _, ok := newProxyRequest(req)
		require.Equal(t, "", hostport)
		require.False(t, ok)
	})

	t.Run("proxy destinations", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/proxy/svc-a.ns1/proxy/svc-b.ns2/proxy/svc-c.ns3", nil)
		hostport, newReq, ok := newProxyRequest(req)
		require.Equal(t, "svc-a.ns1", hostport)
		require.Equal(t, "/proxy/svc-b.ns2/proxy/svc-c.ns3", newReq.URL.String())
		require.True(t, ok)
	})
}
