package http

import (
	"testing"
)

func Test_extractURLForCorsConf(t *testing.T) {
	tests := []struct {
		name string
		urls string
		want string
	}{
		{name: "URL with slash suffix", urls: "https://qwqaq.com/", want: "https://qwqaq.com"},
		{name: "URL with path", urls: "https://qwqaq.com/test-page/", want: "https://qwqaq.com"},
		{name: "URL with port and path", urls: "https://qwqaq.com:12345/test-page/", want: "https://qwqaq.com:12345"},
		{name: "URL with http schema", urls: "http://qwqaq.com", want: "http://qwqaq.com"},
		{name: "URL with regexp, port and path", urls: "http://*.qwqaq.com:12345/test-page/", want: "http://*.qwqaq.com:12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractURLForCorsConf(tt.urls); got != tt.want {
				t.Errorf("extractURLForCorsConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
