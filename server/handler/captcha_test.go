package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func captchaRouter(conf config.CaptchaConf) *fiber.App {
	app := core.NewApp(&config.Config{Captcha: conf})
	router := fiber.New()
	router.Use(limiter.ActionLimitMiddleware(app, limiter.ActionLimitConf{}))
	handler.CaptchaGet(app, router)
	handler.CaptchaVerify(app, router)
	handler.CaptchaStatus(app, router)
	router.Post("/protected", common.LimiterGuard(app, func(c *fiber.Ctx) error { return c.SendStatus(200) }))
	return router
}

func captchaRequest(t *testing.T, router *fiber.App, method, path, body string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := router.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, string(data)
}

func TestCapCaptchaFlow(t *testing.T) {
	used := false
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("response") != "one-use-token" || used {
			w.WriteHeader(404)
			return
		}
		used = true
		_, _ = io.WriteString(w, `{"success":true}`)
	}))
	defer upstream.Close()
	router := captchaRouter(config.CaptchaConf{Enabled: true, Always: true, CaptchaType: config.TypeCapJS, CapJS: config.CapJSConf{APIEndpoint: upstream.URL, KeyID: "site-id", SecretKey: "test-secret"}})
	status := func(want bool) {
		t.Helper()
		code, body := captchaRequest(t, router, "GET", "/captcha/status", "")
		require.Equal(t, 200, code)
		var result handler.ResponseCaptchaStatus
		require.NoError(t, json.Unmarshal([]byte(body), &result))
		require.Equal(t, want, result.IsPass)
	}
	code, page := captchaRequest(t, router, "GET", "/captcha/", "")
	require.Equal(t, 200, code)
	require.Contains(t, page, "cap-widget")
	require.NotContains(t, page, "test-secret")
	status(false)
	code, _ = captchaRequest(t, router, "POST", "/protected", "")
	require.Equal(t, 403, code)
	code, _ = captchaRequest(t, router, "POST", "/captcha/verify", `{"value":"invalid"}`)
	require.Equal(t, 403, code)
	status(false)
	code, _ = captchaRequest(t, router, "POST", "/captcha/verify", `{"value":"one-use-token"}`)
	require.Equal(t, 200, code)
	status(true)
	code, _ = captchaRequest(t, router, "POST", "/protected", "")
	require.Equal(t, 200, code)
	status(false)
	code, _ = captchaRequest(t, router, "POST", "/captcha/verify", `{"value":"one-use-token"}`)
	require.Equal(t, 403, code)
	status(false)
}

type captchaTransport func(*http.Request) (*http.Response, error)

func (f captchaTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCaptchaRejectedWithoutError(t *testing.T) {
	// Existing iframe checkers can return (false, nil) on a non-200 response.
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = captchaTransport(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 503, Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header), Request: r}, nil
	})
	router := captchaRouter(config.CaptchaConf{Enabled: true, Always: true, CaptchaType: config.TypeTurnstile})
	code, body := captchaRequest(t, router, "POST", "/captcha/verify", `{"value":"token"}`)
	require.Equal(t, 403, code)
	require.Contains(t, body, "Verification failed")
	code, _ = captchaRequest(t, router, "POST", "/protected", "")
	require.Equal(t, 403, code)
}
