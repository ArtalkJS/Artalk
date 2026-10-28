package middleware_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/middleware"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestTrustedOriginCors(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()
	app.Conf().TrustedDomains = []string{" https://exact.example.com/blog ", "https://*.example.com:8443", "*"}
	require.NoError(t, app.Dao().DB().Model(&entity.Site{}).Where("id = ?", 1000).
		Update("urls", "https://site.example.org/blog,https://*.example.org/blog,http://localhost:8080,http://[::1]:8080").Error)

	router := fiber.New()
	router.Use(middleware.CorsMiddleware(app.App))
	router.Get("/", func(c *fiber.Ctx) error { return c.SendString("ok") })
	for _, tc := range []struct {
		origin  string
		allowed bool
	}{
		{"https://exact.example.com", true},
		{"https://sub.example.com:8443", true},
		{"https://site.example.org", true},
		{"https://sub.example.org", true},
		{"http://localhost:8080", true},
		{"http://[::1]:8080", true},
		{"https://example.org", false},
		{"https://a.b.example.org", false},
		{"http://sub.example.org", false},
		{"https://sub.example.com", false},
		{"https://sub.example.org.evil", false},
		{"https://untrusted.test", false},
		{"https://evil?sub.example.org", false},
		{"null", false},
	} {
		t.Run(tc.origin, func(t *testing.T) {
			require.Equal(t, tc.allowed, middleware.CheckOriginTrusted(app.App, tc.origin))
			for _, method := range []string{"GET", "OPTIONS"} {
				request := httptest.NewRequest(method, "/", nil)
				request.Header.Set("Origin", tc.origin)
				if method == "OPTIONS" {
					request.Header.Set("Access-Control-Request-Method", "POST")
					request.Header.Set("Access-Control-Request-Headers", "content-type,authorization,x-artalk-guest")
				}
				response, err := router.Test(request)
				require.NoError(t, err)
				response.Body.Close()
				expectedOrigin := ""
				if tc.allowed {
					expectedOrigin = tc.origin
				}
				require.Equal(t, expectedOrigin, response.Header.Get("Access-Control-Allow-Origin"), method)
				if tc.allowed && method == "OPTIONS" {
					require.Equal(t, 204, response.StatusCode)
					require.Contains(t, strings.ToLower(response.Header.Get("Access-Control-Allow-Headers")), "x-artalk-guest")
					require.Contains(t, strings.ToLower(response.Header.Get("Access-Control-Allow-Headers")), "authorization")
				}
			}
		})
	}
	// URL checks still extract the actual host before matching a wildcard.
	trusted, origin, err := middleware.CheckURLTrusted(app.App, "https://sub.example.org/post?x=1")
	require.NoError(t, err)
	require.True(t, trusted)
	require.Equal(t, "https://sub.example.org", origin)
	trusted, _, err = middleware.CheckURLTrusted(app.App, "https://sub.example.org@evil.test/post")
	require.NoError(t, err)
	require.False(t, trusted)
}
