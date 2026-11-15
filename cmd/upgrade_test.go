package cmd

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

type updateTransport func(*http.Request) (*http.Response, error)

func (f updateTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestUpgradeWithoutApplicableRelease(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	for _, tc := range []struct {
		name, body string
		force      bool
	}{
		{"empty release list", "[]", false},
		{"force with no compatible release", "[]", true},
		{"current release", fmt.Sprintf(`[{"tag_name":%q,"assets":[{"id":1,"name":%q,"browser_download_url":"https://fixture.invalid/never-download"}]}]`, config.Version, "artalk_"+runtime.GOOS+"_"+runtime.GOARCH+".tar.gz"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			http.DefaultTransport = updateTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				require.Equal(t, "https://api.github.com/repos/ArtalkJS/Artalk/releases", r.URL.String())
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tc.body)), Request: r}, nil
			})
			command := NewUpgradeCommand()
			if tc.force {
				command.SetArgs([]string{"--force"})
			} else {
				command.SetArgs([]string{})
			}
			require.NoError(t, command.Execute())
			require.Equal(t, 1, calls, "no binary download should occur")
		})
	}
}
