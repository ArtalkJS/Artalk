package artransfer

import "testing"

func Test_stripDomainFromURL(t *testing.T) {
	type args struct {
		fullURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{"https://github.com/abc"}, "/abc"},
		{"2", args{"https://domain.com/some/path?query=123#section"}, "/some/path?query=123#section"},
		{"3", args{"https://domain.com"}, "/"},
		{"4", args{"/some/path?query=123#section"}, "/some/path?query=123#section"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripDomainFromURL(tt.args.fullURL); got != tt.want {
				t.Errorf("stripDomainFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
