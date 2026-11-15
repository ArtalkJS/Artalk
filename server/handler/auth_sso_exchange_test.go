package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthSSOExchangeRequiresVerifiedEmail(t *testing.T) {
	tests := []struct {
		description string
		userinfo    string
	}{
		{
			description: "Reject explicitly unverified admin email",
			userinfo:    `{"sub":"attacker","email":"admin@qwqaq.com","email_verified":false,"nickname":"admin"}`,
		},
		{
			description: "Reject missing email verification claim",
			userinfo:    `{"sub":"attacker","email":"admin@qwqaq.com","nickname":"admin"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/userinfo", r.URL.Path)
				assert.Equal(t, "Bearer attacker-token", r.Header.Get("Authorization"))
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write([]byte(tt.userinfo))
				assert.NoError(t, err)
			}))
			defer issuer.Close()

			app, fiberApp := NewApiTestApp()
			defer app.Cleanup()
			app.Conf().Auth.Enabled = true
			app.Conf().Auth.SSO.Enabled = true
			app.Conf().Auth.SSO.Issuer = issuer.URL
			handler.AuthSSOExchange(app.App, fiberApp)
			assert.True(t, app.Dao().FindUser("admin", "admin@qwqaq.com").IsAdmin)

			req := httptest.NewRequest(http.MethodPost, "/sso/exchange", bytes.NewBufferString(`{"token":"attacker-token"}`))
			req.Header.Set("Content-Type", "application/json")
			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.JSONEq(t, `{"msg":"SSO email is not verified"}`, string(body))
		})
	}
}

func TestAuthSSOExchangeAcceptsVerifiedEmail(t *testing.T) {
	issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/userinfo", r.URL.Path)
		assert.Equal(t, "Bearer verified-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"sub":"verified-user","email":"verified@example.com","email_verified":true,"nickname":"verified"}`))
		assert.NoError(t, err)
	}))
	defer issuer.Close()

	app, fiberApp := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Auth.Enabled = true
	app.Conf().Auth.SSO.Enabled = true
	app.Conf().Auth.SSO.Issuer = issuer.URL
	handler.AuthSSOExchange(app.App, fiberApp)

	req := httptest.NewRequest(http.MethodPost, "/sso/exchange", bytes.NewBufferString(`{"token":"verified-token"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := fiberApp.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result handler.ResponseUserLogin
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "verified", result.User.Name)
	assert.Equal(t, "verified@example.com", result.User.Email)
}

func TestSSOProviderFailuresAndNameFallbacks(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Auth.Enabled = true
	app.Conf().Auth.SSO.Enabled = true
	var status int
	var response string
	issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		io.WriteString(w, response)
	}))
	t.Cleanup(issuer.Close)
	app.Conf().Auth.SSO.Issuer = issuer.URL
	handler.AuthSSOExchange(app.App, router)
	for _, tc := range []struct {
		name, body, wantName string
		providerStatus, want int
	}{
		{"provider rejects token", `{}`, "", 401, 401},
		{"malformed response", `invalid`, "", 200, 502},
		{"missing email", `{"email_verified":true}`, "", 200, 400},
		{"given name", `{"email":"reader@example.org","email_verified":true,"given_name":"Given"}`, "Given", 200, 200},
		{"full name", `{"email":"reader@example.org","email_verified":true,"name":"Full name"}`, "Full name", 200, 200},
		{"email fallback", `{"email":" Reader@Example.org ","email_verified":true}`, "reader@example.org", 200, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			status, response = tc.providerStatus, tc.body
			body := flowRequest(t, router, http.MethodPost, "/sso/exchange", handler.ParamsAuthSSOExchange{Token: "test-token"}, "", tc.want)
			if tc.want == 200 {
				var result handler.ResponseUserLogin
				require.NoError(t, json.Unmarshal(body, &result))
				require.NotEmpty(t, result.Token)
				require.Equal(t, "reader@example.org", result.User.Email)
				require.Equal(t, tc.wantName, result.User.Name)
				require.False(t, result.User.IsAdmin)
			}
		})
	}
	flowRequest(t, router, http.MethodPost, "/sso/exchange", handler.ParamsAuthSSOExchange{}, "", 400)
	app.Conf().Auth.SSO.Issuer = ""
	flowRequest(t, router, http.MethodPost, "/sso/exchange", handler.ParamsAuthSSOExchange{Token: "test"}, "", 500)
	app.Conf().Auth.SSO.Issuer = "http://[invalid"
	flowRequest(t, router, http.MethodPost, "/sso/exchange", handler.ParamsAuthSSOExchange{Token: "test"}, "", 500)
	issuer.Close()
	app.Conf().Auth.SSO.Issuer = issuer.URL
	flowRequest(t, router, http.MethodPost, "/sso/exchange", handler.ParamsAuthSSOExchange{Token: "test"}, "", 503)
	app.Conf().Auth.SSO.Enabled = false
	flowRequest(t, router, http.MethodPost, "/sso/exchange", handler.ParamsAuthSSOExchange{Token: "test"}, "", 404)
}
