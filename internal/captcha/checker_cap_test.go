package captcha

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCapVerify(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		pass   bool
	}{
		{"success", 200, `{"success":true}`, true},
		{"additional fields", 200, `{"success":true,"extra":"ok"}`, true},
		{"rejected", 200, `{"success":false}`, false},
		{"missing result", 200, `{}`, false},
		{"null result", 200, `{"success":null}`, false},
		{"string result", 200, `{"success":"true"}`, false},
		{"number result", 200, `{"success":1}`, false},
		{"invalid JSON", 200, `{"success":true`, false},
		{"extra JSON", 200, `{"success":true}{}`, false},
		{"HTML", 200, `<html>bad gateway</html>`, false},
		{"oversized", 200, `{"success":true,"padding":"` + strings.Repeat("x", 64<<10) + `"}`, false},
		{"invalid token", 400, `{"success":false}`, false},
		{"bad secret", 403, `{"success":false}`, false},
		{"used token", 404, `{"success":false}`, false},
		{"rate limited", 429, `{"success":false}`, false},
		{"upstream error", 500, `{"success":true}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/cap/site-id/siteverify", r.URL.Path)
				require.Equal(t, "POST", r.Method)
				require.NoError(t, r.ParseForm())
				require.Equal(t, "test-secret", r.Form.Get("secret"))
				require.Equal(t, "test-token", r.Form.Get("response"))
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer upstream.Close()
			checker := NewCapJSChecker(&config.CapJSConf{APIEndpoint: upstream.URL + "/cap/", KeyID: "site-id", SecretKey: "test-secret"}, &User{})
			pass, err := checker.Check("test-token")
			require.Equal(t, tc.pass, pass)
			if tc.pass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.NotContains(t, err.Error(), "test-secret")
			}
		})
	}
}

type capTransport func(*http.Request) (*http.Response, error)

func (f capTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type capBody struct {
	io.Reader
	closed bool
}

func (b *capBody) Close() error { b.closed = true; return nil }

type truncatedCapBody struct{}

func (truncatedCapBody) Read(p []byte) (int, error) {
	return copy(p, `{"success":true}`), io.ErrUnexpectedEOF
}

func TestCapTransportFailures(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	checker := NewCapJSChecker(&config.CapJSConf{APIEndpoint: "https://cap.example", KeyID: "site-id", SecretKey: "test-secret"}, &User{})
	for _, status := range []int{200, 400, 403, 404, 429, 500} {
		body := &capBody{Reader: strings.NewReader(`{"success":true}`)}
		http.DefaultTransport = capTransport(func(r *http.Request) (*http.Response, error) {
			_, hasDeadline := r.Context().Deadline()
			require.True(t, hasDeadline)
			return &http.Response{StatusCode: status, Body: body, Header: make(http.Header), Request: r}, nil
		})
		_, _ = checker.Check("token")
		require.True(t, body.closed, "HTTP %d body must be closed", status)
	}
	body := &capBody{Reader: truncatedCapBody{}}
	http.DefaultTransport = capTransport(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: body, Header: make(http.Header), Request: r}, nil
	})
	pass, err := checker.Check("token")
	require.False(t, pass)
	require.Error(t, err)
	require.True(t, body.closed)
	http.DefaultTransport = capTransport(func(*http.Request) (*http.Response, error) { return nil, context.DeadlineExceeded })
	pass, err = checker.Check("token")
	require.False(t, pass)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestCapRejectsRedirects(t *testing.T) {
	hit := false
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hit = true; _, _ = io.WriteString(w, `{"success":true}`) }))
	defer destination.Close()
	for _, status := range []int{301, 302, 307, 308} {
		source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, destination.URL, status) }))
		checker := NewCapJSChecker(&config.CapJSConf{APIEndpoint: source.URL, KeyID: "site-id", SecretKey: "test-secret"}, &User{})
		pass, err := checker.Check("token")
		require.False(t, pass)
		require.Error(t, err)
		source.Close()
	}
	require.False(t, hit, "do not forward verification credentials through a redirect")
}

func TestCapConfigurationAndPage(t *testing.T) {
	conf := config.CapJSConf{APIEndpoint: " https://cap.example/subpath/ ", KeyID: "site-id", SecretKey: "test-secret"}
	checker := NewCapJSChecker(&conf, &User{})
	page, err := checker.Get()
	require.NoError(t, err)
	require.NotContains(t, string(page), conf.SecretKey)
	require.Contains(t, string(page), `https://cap.example/subpath/site-id/`)
	require.Contains(t, string(page), `https:\/\/cap.example\/subpath/assets/hashwx.wasm`)
	require.Contains(t, string(page), `https:\/\/cap.example\/subpath/assets/cap_wasm_bg.wasm`)
	require.NotContains(t, string(page), "console.log")
	require.Equal(t, IFrame, checker.Type())
	for _, endpoint := range []string{"", "cap.example", "ftp://cap.example", "https://user:password@cap.example", "https://cap.example?x=1", "https://cap.example?", "https://cap.example#fragment", "https://bad%host"} {
		invalid := conf
		invalid.APIEndpoint = endpoint
		_, err := NewCapJSChecker(&invalid, &User{}).Get()
		require.Error(t, err, endpoint)
	}
	for _, key := range []string{"", "..", ".", "x/y", "x?y", "x#y", "x%2fy"} {
		invalid := conf
		invalid.KeyID = key
		_, err := NewCapJSChecker(&invalid, &User{}).Get()
		require.Error(t, err, key)
	}
	invalid := conf
	invalid.SecretKey = ""
	_, err = NewCapJSChecker(&invalid, &User{}).Get()
	require.Error(t, err)
	pass, err := checker.Check("")
	require.False(t, pass)
	require.Error(t, err)
}
