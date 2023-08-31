package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarked(t *testing.T) {
	html, err := Marked("# Title")
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
}
