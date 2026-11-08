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
		{"match another subdomain", "https://*.example.org", "https://blog.example.org", true},
		{"match hyphen", "https://*.example.com", "https://my-blog.example.com", true},

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

		// Parsed hosts, ports, and configured URL paths.
		{"explicit port", "https://*.example.com:8443", "https://sub.example.com:8443", true},
		{"different port", "https://*.example.com:8443", "https://sub.example.com:9443", false},
		{"missing port", "https://*.example.com:8443", "https://sub.example.com", false},
		{"extra port", "https://*.example.com", "https://sub.example.com:8443", false},
		{"configured path", " https://*.example.com/blog/ ", "https://sub.example.com", true},
		{"host case", "https://*.EXAMPLE.com", "https://Sub.example.COM", true},
		{"punycode label", "https://*.example.com", "https://xn--bcher-kva.example.com", true},

		// A source must be an origin, not a URL containing credentials or a path.
		{"path suffix", "https://*.example.com", "https://evil/@sub.example.com", false},
		{"query suffix", "https://*.example.com", "https://evil?sub.example.com", false},
		{"fragment suffix", "https://*.example.com", "https://evil#sub.example.com", false},
		{"credentials", "https://*.example.com", "https://user@sub.example.com", false},
		{"origin path", "https://*.example.com", "https://sub.example.com/", false},
		{"empty query", "https://*.example.com", "https://sub.example.com?", false},
		{"empty fragment", "https://*.example.com", "https://sub.example.com#", false},
		{"encoded separator", "https://*.example.com", "https://evil%2fsub.example.com", false},
		{"domain suffix attack", "https://*.example.com", "https://sub.example.com.evil", false},
		{"leading hyphen", "https://*.example.com", "https://-sub.example.com", false},
		{"trailing hyphen", "https://*.example.com", "https://sub-.example.com", false},
		{"literal wildcard", "https://*.example.com", "https://*.example.com", false},
		{"null origin", "https://*.example.com", "null", false},
		{"empty origin port", "https://*.example.com", "https://sub.example.com:", false},
		{"empty pattern port", "https://*.example.com:", "https://sub.example.com", false},
		{"empty base domain", "https://*.", "https://sub.", false},
		{"invalid port", "https://*.example.com", "https://sub.example.com:bad", false},
		{"bare wildcard", "*", "https://sub.example.com", false},
		{"multiple wildcards", "https://*.*.example.com", "https://a.b.example.com", false},
		{"partial wildcard", "https://foo*.example.com", "https://foobar.example.com", false},
		{"pattern credentials", "https://user@*.example.com", "https://sub.example.com", false},
		{"invalid pattern", "https://%zz.example.com", "https://sub.example.com", false},
		{"non HTTP scheme", "ftp://*.example.com", "ftp://sub.example.com", false},

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
