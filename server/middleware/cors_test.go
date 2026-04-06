package middleware

import (
	"testing"
)

func TestMatchWildcardOrigin(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		origin  string
		want    bool
	}{
		// Basic wildcard matching
		{"match subdomain", "https://*.example.com", "https://sub.example.com", true},
		{"match another subdomain", "https://*.wgfilm21.net", "https://go6.wgfilm21.net", true},
		{"match alternatif", "https://*.wgfilm21.net", "https://alternatif.wgfilm21.net", true},

		// Should not match
		{"no match bare domain", "https://*.example.com", "https://example.com", false},
		{"no match nested subdomain", "https://*.example.com", "https://a.b.example.com", false},
		{"no match different domain", "https://*.example.com", "https://sub.other.com", false},
		{"no match different scheme", "https://*.example.com", "http://sub.example.com", false},

		// Non-wildcard patterns should not match
		{"non-wildcard exact", "https://example.com", "https://example.com", false},
		{"non-wildcard no match", "https://example.com", "https://sub.example.com", false},

		// Edge cases
		{"empty origin", "https://*.example.com", "", false},
		{"empty pattern", "", "https://sub.example.com", false},
		{"no scheme in pattern", "*.example.com", "https://sub.example.com", false},
		{"no scheme in origin", "https://*.example.com", "sub.example.com", false},

		// HTTP scheme
		{"http wildcard", "http://*.example.com", "http://sub.example.com", true},
		{"http no match https", "http://*.example.com", "https://sub.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchWildcardOrigin(tt.pattern, tt.origin)
			if got != tt.want {
				t.Errorf("matchWildcardOrigin(%q, %q) = %v, want %v", tt.pattern, tt.origin, got, tt.want)
			}
		})
	}
}
