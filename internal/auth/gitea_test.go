package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
)

func TestGiteaSiteURL(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"", "https://gitea.com"},
		{" https://GITEA.com:443/ ", "https://gitea.com"},
		{"http://localhost:8080/git/", "http://localhost:8080/git"},
		{"http://[::1]:8080/git/", "http://[::1]:8080/git"},
		{"http://git.example:80/", "http://git.example"},
		{"https://git.example/team%20one/", "https://git.example/team%20one"},
	} {
		got, err := giteaSiteURL(tc.input)
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
	for _, raw := range []string{"git.example", "//git.example", "ftp://git.example", "https://", "https://%", "https://user:secret@git.example", "https://git.example/?token=x", "https://git.example/?", "https://git.example/#fragment"} {
		_, err := giteaSiteURL(raw)
		require.Error(t, err, raw)
		conf := &config.Config{}
		conf.Auth.Gitea.Enabled = true
		conf.Auth.Gitea.SiteURL = raw
		require.Empty(t, GetProviders(conf), "invalid URL must not fall back to gitea.com")
	}
}

func TestGiteaOAuthWithSubpath(t *testing.T) {
	var requested []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = append(requested, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/git/login/oauth/access_token":
			fmt.Fprint(w, `{"access_token":"test-token","token_type":"Bearer"}`)
		case "/git/api/v1/user":
			fmt.Fprint(w, `{"id":42,"login":"alice","full_name":"Alice","email":"alice@example.test"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	conf := &config.Config{}
	conf.Auth.Callback = "https://comments.example/api/v2/auth/{provider}/callback"
	conf.Auth.Gitea.Enabled = true
	conf.Auth.Gitea.SiteURL = server.URL + "/git/"
	conf.Auth.Gitea.ClientID = "client"
	conf.Auth.Gitea.ClientSecret = "secret"
	providers := GetProviders(conf)
	require.Len(t, providers, 1)
	provider := providers[0]
	require.Equal(t, "gitea", provider.Name())
	session, err := provider.BeginAuth("test-state")
	require.NoError(t, err)
	authURL, err := session.GetAuthURL()
	require.NoError(t, err)
	u, err := url.Parse(authURL)
	require.NoError(t, err)
	require.Equal(t, "/git/login/oauth/authorize", u.Path)
	require.Equal(t, "https://comments.example/api/v2/auth/gitea/callback", u.Query().Get("redirect_uri"))
	require.Equal(t, "test-state", u.Query().Get("state"))
	_, err = session.Authorize(provider, url.Values{"code": {"test-code"}})
	require.NoError(t, err)
	user, err := provider.FetchUser(session)
	require.NoError(t, err)
	require.Equal(t, "42", user.UserID)
	require.Equal(t, []string{"/git/login/oauth/access_token", "/git/api/v1/user"}, requested)
	social, err := GetSocialUser(user, conf.Auth.Gitea.SiteURL)
	require.NoError(t, err)
	require.NotEqual(t, "42", social.RemoteUID)
	conf.Auth.Gitea.Enabled = false
	require.Empty(t, GetProviders(conf))
}

func TestGiteaIdentityCompatibility(t *testing.T) {
	user := goth.User{Provider: "gitea", UserID: "42", Name: "Alice", Email: "alice@example.test"}
	identities := map[string]string{}
	for _, site := range []string{"", "https://gitea.com/", "https://GITEA.COM:443", "https://git.example/a", "https://git.example/b", "https://other.example/a"} {
		conf := &config.Config{}
		conf.Auth.Gitea.SiteURL = site
		social, err := GetSocialUser(user, conf.Auth.Gitea.SiteURL)
		require.NoError(t, err)
		identities[site] = social.RemoteUID
		require.LessOrEqual(t, len(social.RemoteUID), 255)
	}
	require.Equal(t, "42", identities[""])
	require.Equal(t, "42", identities["https://gitea.com/"])
	require.Equal(t, "42", identities["https://GITEA.COM:443"])
	require.NotEqual(t, identities[""], identities["https://git.example/a"])
	require.NotEqual(t, identities["https://git.example/a"], identities["https://git.example/b"])
	require.NotEqual(t, identities["https://git.example/a"], identities["https://other.example/a"])
	conf := &config.Config{}
	for _, id := range []string{"", "0"} {
		user.UserID = id
		_, err := GetSocialUser(user, conf.Auth.Gitea.SiteURL)
		require.Error(t, err)
	}
}
