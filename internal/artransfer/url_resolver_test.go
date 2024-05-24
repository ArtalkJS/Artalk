package artransfer

import (
	"testing"
)

func Test_getResolvedPageKey(t *testing.T) {
	type args struct {
		baseUrlRaw    string
		commentUrlRaw string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Domain-only combine with Path-only", args{"https://github.com", "/1.html"}, "https://github.com/1.html"},
		{"Domain-only combine with Domain-Path", args{"https://github.com", "https://xxx.com/1.html"}, "https://github.com/1.html"},
		{"Domain-only (trailing slash) combine with Path-only", args{"https://github.com/", "/1.html"}, "https://github.com/1.html"},
		{"Empty combine with Path-only", args{"", "/1.html"}, "/1.html"},
		{"Empty combine with Domain-Path", args{"", "https://xxx.com/1.html"}, "https://xxx.com/1.html"},
		{"Domain-Path (no trailing slash) combine with Path-only (trailing slash)", args{"https://github.com/233", "/1/"}, "https://github.com/1/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResolvedPageKey(tt.args.baseUrlRaw, tt.args.commentUrlRaw); got != tt.want {
				t.Errorf("getResolvedPageKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
