package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/okta"
	"github.com/stretchr/testify/require"
)

func TestOktaIssuer(t *testing.T) {
	for _, tc := range []struct{ raw, want string }{
		{" https://TENANT.okta.com:443/ ", "https://tenant.okta.com"},
		{"https://tenant.okta.com/oauth2/default/", "https://tenant.okta.com/oauth2/default"},
		{"https://login.example:8443/oauth2/custom", "https://login.example:8443/oauth2/custom"},
	} {
		got, err := oktaIssuerURL(tc.raw)
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
	for _, raw := range []string{"", "tenant.okta.com", "//tenant.okta.com", "http://tenant.okta.com", "ftp://tenant.okta.com", "https://", "https://%", "https://user:secret@tenant.okta.com", "https://tenant.okta.com?x=1", "https://tenant.okta.com?", "https://tenant.okta.com/#fragment"} {
		conf := &config.Config{}
		conf.Auth.Okta.Enabled = true
		conf.Auth.Okta.Issuer = raw
		_, err := oktaIssuerURL(raw)
		require.Error(t, err, raw)
		require.Empty(t, GetProviders(conf), raw)
	}
}

func TestOktaOAuth(t *testing.T) {
	for _, path := range []string{"", "/oauth2/default", "/oauth2/custom"} {
		t.Run(path, func(t *testing.T) {
			basePath := path
			if basePath == "" {
				basePath = "/oauth2"
			}
			var requested []string
			upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requested = append(requested, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case basePath + "/v1/token":
					require.Equal(t, "POST", r.Method)
					require.Equal(t, "test-code", r.FormValue("code"))
					require.Equal(t, "https://comments.example/api/v2/auth/okta/callback", r.FormValue("redirect_uri"))
					client, secret, ok := r.BasicAuth()
					require.True(t, ok)
					require.Equal(t, "client", client)
					require.Equal(t, "secret", secret)
					fmt.Fprint(w, `{"access_token":"local-token","token_type":"Bearer"}`)
				case basePath + "/v1/userinfo":
					require.Equal(t, "Bearer local-token", r.Header.Get("Authorization"))
					fmt.Fprint(w, `{"sub":"user-42","name":"Alice","email":"alice@example.test","email_verified":false}`)
				default:
					w.WriteHeader(404)
				}
			}))
			defer upstream.Close()
			conf := &config.Config{}
			conf.Auth.Callback = "https://comments.example/api/v2/auth/{provider}/callback"
			conf.Auth.Okta.Enabled = true
			conf.Auth.Okta.Issuer = upstream.URL + path + "/"
			conf.Auth.Okta.ClientID = "client"
			conf.Auth.Okta.ClientSecret = "secret"
			providers := GetProviders(conf)
			require.Len(t, providers, 1)
			provider := providers[0].(*okta.Provider)
			provider.HTTPClient.Transport = upstream.Client().Transport
			session, err := provider.BeginAuth("state")
			require.NoError(t, err)
			raw, err := session.GetAuthURL()
			require.NoError(t, err)
			u, err := url.Parse(raw)
			require.NoError(t, err)
			require.Equal(t, basePath+"/v1/authorize", u.Path)
			require.Equal(t, "state", u.Query().Get("state"))
			require.Equal(t, "https://comments.example/api/v2/auth/okta/callback", u.Query().Get("redirect_uri"))
			require.ElementsMatch(t, []string{"openid", "profile", "email"}, strings.Fields(u.Query().Get("scope")))
			_, err = session.Authorize(provider, url.Values{"code": {"test-code"}})
			require.NoError(t, err)
			profile, err := provider.FetchUser(session)
			require.NoError(t, err)
			require.Equal(t, "user-42", profile.UserID)
			social, err := GetSocialUser(profile, &conf.Auth)
			require.NoError(t, err)
			require.NotEqual(t, profile.UserID, social.RemoteUID)
			require.Equal(t, []string{basePath + "/v1/token", basePath + "/v1/userinfo"}, requested)
			require.NotEmpty(t, GetProviderInfo(conf, providers)[0].Icon)
			conf.Auth.Okta.Enabled = false
			require.Empty(t, GetProviders(conf))
		})
	}
}

func TestOktaIdentity(t *testing.T) {
	conf := &config.AuthConf{}
	user := goth.User{Provider: "okta", UserID: "same-subject", Name: "Alice", Email: "alice@example.test"}
	ids := map[string]string{}
	for _, issuer := range []string{"https://tenant.okta.com", "https://TENANT.okta.com:443/", "https://tenant.okta.com/oauth2/default", "https://tenant.okta.com/oauth2/custom", "https://other.okta.com"} {
		conf.Okta.Issuer = issuer
		social, err := GetSocialUser(user, conf)
		require.NoError(t, err)
		ids[issuer] = social.RemoteUID
	}
	require.Equal(t, ids["https://tenant.okta.com"], ids["https://TENANT.okta.com:443/"])
	require.NotEqual(t, ids["https://tenant.okta.com"], ids["https://tenant.okta.com/oauth2/default"])
	require.NotEqual(t, ids["https://tenant.okta.com/oauth2/default"], ids["https://tenant.okta.com/oauth2/custom"])
	require.NotEqual(t, ids["https://tenant.okta.com"], ids["https://other.okta.com"])
	for _, id := range []string{"", "   "} {
		user.UserID = id
		_, err := GetSocialUser(user, conf)
		require.Error(t, err)
	}
	user.UserID = strings.Repeat("a", 255)
	social, err := GetSocialUser(user, conf)
	require.NoError(t, err)
	require.LessOrEqual(t, len(social.RemoteUID), 255)
	user.UserID = "another-subject"
	other, err := GetSocialUser(user, conf)
	require.NoError(t, err)
	require.NotEqual(t, social.RemoteUID, other.RemoteUID)
	// Other providers retain their existing identity and profile fallback rules.
	user = goth.User{Provider: "github", UserID: "42", Email: "alice@example.test"}
	social, err = GetSocialUser(user, conf)
	require.NoError(t, err)
	require.Equal(t, "42", social.RemoteUID)
	require.Equal(t, "alice", social.Name)
}

func TestOktaHTTPFailures(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusTemporaryRedirect} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			reached := false
			destination := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached = true }))
			defer destination.Close()
			upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", destination.URL)
				w.WriteHeader(status)
				io.WriteString(w, `{"access_token":"bad-token","sub":"bad-user"}`)
			}))
			defer upstream.Close()
			provider, err := newOktaProvider(upstream.URL, "client", "secret", "https://comments.example/callback")
			require.NoError(t, err)
			provider.HTTPClient.Transport = upstream.Client().Transport
			session, err := provider.BeginAuth("state")
			require.NoError(t, err)
			_, err = session.Authorize(provider, url.Values{"code": {"code"}})
			require.Error(t, err)
			_, err = provider.FetchUser(&okta.Session{AccessToken: "local-token"})
			require.Error(t, err)
			require.False(t, reached, "credentials and access tokens must not follow redirects")
		})
	}
}
