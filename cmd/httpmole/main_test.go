package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToHeadersMap(t *testing.T) {
	h := toHeadersMap([]string{"a:b"})
	require.Equal(t, "b", h.Get("a"))
}
