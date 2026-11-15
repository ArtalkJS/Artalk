package gothic_fiber

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/stretchr/testify/require"
)

func TestOAuthSessionLifecycle(t *testing.T) {
	previousStore, previousProviders := SessionStore, goth.GetProviders()
	t.Cleanup(func() {
		SessionStore = previousStore
		goth.ClearProviders()
		for _, p := range previousProviders {
			goth.UseProviders(p)
		}
	})
	SessionStore = session.New()
	goth.ClearProviders()
	goth.UseProviders(github.New("client", "secret", "https://example.org/callback"))
	router := fiber.New()
	router.Get("/begin/:provider?", BeginAuthHandler)
	router.Get("/provider/:provider?", func(c *fiber.Ctx) error {
		p, err := GetProviderName(c)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.SendString(p)
	})
	router.Get("/store", func(c *fiber.Ctx) error { return StoreInSession("github", c.Query("value"), c) })
	router.Get("/value", func(c *fiber.Ctx) error {
		value, err := GetFromSession("github", c)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.SendString(value)
	})
	router.Get("/complete/:provider?", func(c *fiber.Ctx) error {
		_, err := CompleteUserAuth(c)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.SendStatus(200)
	})
	router.Get("/logout", Logout)
	var cookie *http.Cookie
	request := func(path, provider string, status int) string {
		t.Helper()
		req := httptest.NewRequest("GET", path, nil)
		if cookie != nil {
			req.AddCookie(cookie)
		}
		if provider != "" {
			req.Header.Set("provider", provider)
		}
		resp, err := router.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, status, resp.StatusCode, string(body))
		if cookies := resp.Cookies(); len(cookies) > 0 {
			cookie = cookies[0]
		}
		return string(body)
	}
	require.Contains(t, request("/begin/unknown", "", 400), "provider")
	request("/complete/unknown", "", 400)
	request("/complete", "", 400)
	request("/complete/github", "", 400)
	request("/value", "", 400)
	require.Equal(t, "github", request("/provider?provider=github", "", 200))
	require.Equal(t, "github", request("/provider/github", "", 200))
	require.Equal(t, "github", request("/provider", "github", 200))
	request("/store?value=round-trip", "", 200)
	require.Equal(t, "round-trip", request("/value", "", 200))
	require.Equal(t, "github", request("/provider", "", 200))
	// Malformed provider state must fail and invalidate the session.
	request("/complete/github", "", 400)
	request("/value", "", 400)
	request("/begin/github?state=expected", "", 307)
	require.Contains(t, request("/complete/github?state=wrong", "", 400), "state token mismatch")
	request("/logout", "", 200)
	request("/value", "", 400)
	SessionStore = nil
	request("/begin/github", "", 400)
	request("/complete/github", "", 400)
}

type failedSessionStorage struct{ fiber.Storage }

func (failedSessionStorage) Get(string) ([]byte, error) { return nil, nil }
func (failedSessionStorage) Set(string, []byte, time.Duration) error {
	return errors.New("session storage unavailable")
}

func TestOAuthSessionSaveFailure(t *testing.T) {
	previous := SessionStore
	t.Cleanup(func() { SessionStore = previous })
	SessionStore = session.New(session.Config{Storage: failedSessionStorage{}})
	router := fiber.New()
	router.Get("/", func(c *fiber.Ctx) error { return StoreInSession("provider", "oauth secret", c) })
	resp, err := router.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 500, resp.StatusCode, "login must not continue when its state was not stored")
}

func TestOAuthRejectsCorruptSessionData(t *testing.T) {
	router := fiber.New()
	router.Get("/", func(c *fiber.Ctx) error {
		store, err := SessionStore.Get(c)
		require.NoError(t, err)
		store.Set("corrupt", "not compressed")
		_, err = getSessionValue(store, "corrupt")
		require.Error(t, err)
		require.NoError(t, updateSessionValue(store, "truncated", "secret"))
		data := store.Get("truncated").(string)
		store.Set("truncated", data[:len(data)-5])
		_, err = getSessionValue(store, "truncated")
		require.Error(t, err)
		for _, url := range []string{"", ":invalid", "https://example.org?state=expected"} {
			err := validateState(c, &github.Session{AuthURL: url})
			require.Error(t, err)
		}
		require.Equal(t, "", GetState(c))
		return c.SendString("rejected")
	})
	resp, err := router.Test(httptest.NewRequest("GET", "/", strings.NewReader("")))
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
}
