package utils

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	testCases := []struct {
		email      string
		shouldPass bool
	}{
		{"test@example.com", true},
		{"invalid.email", false},
		{"test.name@example.com", true},
	}

	for _, tc := range testCases {
		result := ValidateEmail(tc.email)
		if result != tc.shouldPass {
			t.Errorf("For email %s, expected pass: %v, but got: %v", tc.email, tc.shouldPass, result)
		}
	}
}

func TestValidateURL(t *testing.T) {
	testCases := []struct {
		url        string
		shouldPass bool
	}{
		{"https://example.com", true},
		{"invalid.url", false},
		{"http://example.com/path?param=value", true},
		{"ftp://example.com", false},
	}

	for _, tc := range testCases {
		result := ValidateURL(tc.url)
		if result != tc.shouldPass {
			t.Errorf("For URL %s, expected pass: %v, but got: %v", tc.url, tc.shouldPass, result)
		}
	}
}
