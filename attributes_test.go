package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttributes_Get(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var attrs Attributes
		gotVal, gotOK := attrs.Get("key")
		assert.Empty(t, gotVal)
		assert.False(t, gotOK)
	})

	t.Run("key not found", func(t *testing.T) {
		attrs := Attributes{}
		gotVal, gotOK := attrs.Get("key")
		assert.Empty(t, gotVal)
		assert.False(t, gotOK)
	})

	t.Run("nil value slice", func(t *testing.T) {
		attrs := Attributes{"key": nil}
		gotVal, gotOK := attrs.Get("key")
		assert.Empty(t, gotVal)
		assert.False(t, gotOK)
	})

	t.Run("empty value slice", func(t *testing.T) {
		attrs := Attributes{"key": []string{}}
		gotVal, gotOK := attrs.Get("key")
		assert.Empty(t, gotVal)
		assert.False(t, gotOK)
	})

	t.Run("empty string", func(t *testing.T) {
		attrs := Attributes{"key": []string{""}}
		gotVal, gotOK := attrs.Get("key")
		assert.Empty(t, gotVal)
		assert.True(t, gotOK)
	})

	t.Run("multiple values", func(t *testing.T) {
		attrs := Attributes{"key": []string{"first", "second"}}
		gotVal, gotOK := attrs.Get("key")
		assert.Equal(t, "second", gotVal)
		assert.True(t, gotOK)
	})
}
