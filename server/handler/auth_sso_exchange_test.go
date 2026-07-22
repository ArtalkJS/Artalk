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
