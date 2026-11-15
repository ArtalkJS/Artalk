package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/auth/gothic_fiber"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
)

func TestSocialLoginWithLocalProvider(t *testing.T) {
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/login/oauth/access_token":
			_, _ = io.WriteString(w, `{"access_token":"local-token","token_type":"Bearer"}`)
		case "/api/v1/user":
			_, _ = io.WriteString(w, `{"id":42,"login":"reader","full_name":"Reader","email":"social@example.org"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(provider.Close)

	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	previousProviders := goth.GetProviders()
	t.Cleanup(func() {
		goth.ClearProviders()
		for _, provider := range previousProviders {
			goth.UseProviders(provider)
		}
	})
	conf := app.Conf()
	conf.Auth.Enabled = true
	conf.LoginTimeout = 3600
	conf.Auth.Callback = "http://comments.example.org/api/v2/auth/{provider}/callback"
	conf.Auth.Gitea.Enabled = true
	conf.Auth.Gitea.SiteURL = provider.URL
	conf.Auth.Gitea.ClientID = "local-client"
	conf.Auth.Gitea.ClientSecret = "local-secret"
	handler.AuthSocialLogin(app.App, router)
	handler.UserInfo(app.App, router)

	var providers handler.ResponseConfAuthProviders
	body := flowRequest(t, router, http.MethodGet, "/conf/auth/providers", nil, "", http.StatusOK)
	require.NoError(t, json.Unmarshal(body, &providers))
	require.Len(t, providers.Providers, 1)
	require.Equal(t, "gitea", providers.Providers[0].Name)
	require.Equal(t, "http://comments.example.org", providers.CallbackOrigin)

	var userID uint
	for _, action := range []string{"first login", "returning login"} {
		t.Run(action, func(t *testing.T) {
			begin, err := router.Test(httptest.NewRequest(http.MethodGet, "/auth/gitea", nil))
			require.NoError(t, err)
			defer begin.Body.Close()
			require.Equal(t, http.StatusTemporaryRedirect, begin.StatusCode)
			target, err := url.Parse(begin.Header.Get("Location"))
			require.NoError(t, err)
			require.Equal(t, provider.URL+"/login/oauth/authorize", target.Scheme+"://"+target.Host+target.Path)
			state := target.Query().Get("state")
			require.NotEmpty(t, state)
			cookies := begin.Cookies()
			require.Len(t, cookies, 1)
			require.Equal(t, gothic_fiber.SessionName, cookies[0].Name)
			require.True(t, cookies[0].HttpOnly)

			path := "/auth/gitea/callback?" + url.Values{"code": {"local-code"}, "state": {state}}.Encode()
			callback := httptest.NewRequest(http.MethodGet, path, nil)
			callback.AddCookie(cookies[0])
			resp, err := router.Test(callback)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Contains(t, resp.Header.Get("Cache-Control"), "no-store")
			page, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			token := regexp.MustCompile(`payload: '([A-Za-z0-9._-]+)'`).FindSubmatch(page)
			require.Len(t, token, 2)

			var info handler.ResponseUserInfo
			body := flowRequest(t, router, http.MethodGet, "/user", nil, string(token[1]), http.StatusOK)
			require.NoError(t, json.Unmarshal(body, &info))
			require.True(t, info.IsLogin)
			require.NotNil(t, info.User)
			require.False(t, info.User.IsAdmin)
			if userID == 0 {
				userID = info.User.ID
			}
			require.NotZero(t, userID)
			require.Equal(t, userID, info.User.ID)
		})
	}

	conf.Auth.Enabled = false
	flowRequest(t, router, http.MethodGet, "/conf/auth/providers", nil, "", http.StatusNotFound)
}
