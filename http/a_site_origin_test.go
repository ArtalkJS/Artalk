package http

import (
	"testing"
)

func Test_GetIsAllowOrigin(t *testing.T) {
	tests := []struct {
		name      string
		origin    string
		allowURLs []string
		want      bool
	}{
		{name: "matched allowURLs with slash suffix", origin: "https://qwqaq.com", allowURLs: []string{"https://qwqaq.com/"}, want: true},
		{name: "matched allowURLs with path", origin: "https://qwqaq.com", allowURLs: []string{"https://qwqaq.com/test-page/"}, want: true},
		{name: "matched allowURLs with port and path", origin: "https://qwqaq.com:12345", allowURLs: []string{"https://qwqaq.com:12345/test-page/"}, want: true},
		{name: "matched allowURLs with http schema", origin: "http://qwqaq.com", allowURLs: []string{"http://qwqaq.com"}, want: true},

		{name: "not matched, port not same", origin: "https://qwqaq.com:1234", allowURLs: []string{"https://qwqaq.com"}, want: false},
		{name: "not matched, protocol not same", origin: "http://qwqaq.com", allowURLs: []string{"https://qwqaq.com"}, want: false},
		{name: "not matched, hostname not same", origin: "https://abc.qwqaq.com", allowURLs: []string{"https://qwqaq.com"}, want: false},

		{name: "invalid origin 1", origin: "qwqaq.com", allowURLs: []string{"https://qwqaq.com"}, want: false},
		{name: "invalid origin 2", origin: "", allowURLs: []string{"https://qwqaq.com"}, want: false},
		{name: "invalid origin 3", origin: "null", allowURLs: []string{"https://qwqaq.com"}, want: false},

		{name: "matched multi-allowUrls", origin: "https://abc.qwqaq.com", allowURLs: []string{"https://aaaa.com", "https://bbb.com", "https://abc.qwqaq.com/abcd"}, want: true},
		{name: "not matched multi-allowUrls", origin: "https://def.qwqaq.com", allowURLs: []string{"https://aaaa.com", "https://bbb.com", "https://abc.qwqaq.com/abcd"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIsAllowOrigin(tt.origin, tt.allowURLs); got != tt.want {
				t.Errorf("GetIsAllowOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}
