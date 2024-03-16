package ui

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed public pages
var public embed.FS

func TestNew(t *testing.T) {
	mux, err := New(Registry{}, public, Options{})
	require.NoError(t, err)
	_ = mux
}
