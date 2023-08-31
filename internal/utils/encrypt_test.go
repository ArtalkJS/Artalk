package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMD5Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"world", "7d793037a0760186574b0282f2f435e7"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
	}

	for _, test := range tests {
		result := GetMD5Hash(test.input)
		assert.Equal(t, test.expected, result)
	}
}
