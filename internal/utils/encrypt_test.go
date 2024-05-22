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

func TestGetSha256Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"world", "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}

	for _, test := range tests {
		result := GetSha256Hash(test.input)
		assert.Equal(t, test.expected, result)
	}
}
