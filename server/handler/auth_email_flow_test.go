package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
)

func TestEmailRegistrationAndPasswordReset(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	app.Conf().LoginTimeout = 3600
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.AuthEmailRegister(app.App, router)
	handler.UserLogin(app.App, router)
	handler.UserInfo(app.App, router)

	const email = "register@example.org"
	register := handler.RequestAuthEmailRegister{
		Name: "New reader", Email: email, Password: "initial-password", Link: "https://reader.example.org",
	}
	challenge := func(code string) {
		t.Helper()
		require.NoError(t, app.Dao().DB().Create(&entity.UserEmailVerify{
			Email: email, Code: code, ExpiresAt: time.Now().Add(5 * time.Minute),
		}).Error)
		register.Code = code
	}
	challenge("721350")
	var registered handler.ResponseUserLogin
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/auth/email/register", register, "", http.StatusOK), &registered))
	require.NotZero(t, registered.User.ID)
	require.Equal(t, register.Name, registered.User.Name)
	require.Equal(t, register.Link, registered.User.Link)

	var info handler.ResponseUserInfo
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/user", nil, registered.Token, http.StatusOK), &info))
	require.True(t, info.IsLogin)
	require.NotNil(t, info.User)
	require.Equal(t, registered.User.ID, info.User.ID)

	challenge("503720")
	oldPassword := register.Password
	register.Name, register.Password = "", "reset-password"
	var reset handler.ResponseUserLogin
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/auth/email/register", register, "", http.StatusOK), &reset))
	require.Equal(t, registered.User.ID, reset.User.ID)

	login := handler.ParamsUserLogin{Email: email, Password: oldPassword}
	flowRequest(t, router, http.MethodPost, "/user/access_token", login, "", http.StatusUnauthorized)
	login.Password = register.Password
	var signedIn handler.ResponseUserLogin
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/user/access_token", login, "", http.StatusOK), &signedIn))
	require.Equal(t, registered.User.ID, signedIn.User.ID)
}

func TestEmailRegistrationValidationPreservesChallenge(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Auth.Email.Enabled = true
	app.Conf().Captcha.Enabled = false
	handler.AuthEmailRegister(app.App, router)
	handler.AuthEmailLogin(app.App, router)
	challenge := entity.UserEmailVerify{Email: "validation@example.org", Code: "654321", ExpiresAt: time.Now().Add(time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	for _, tc := range []struct {
		name string
		edit func(*handler.RequestAuthEmailRegister)
	}{
		{"bad email", func(p *handler.RequestAuthEmailRegister) { p.Email = "invalid" }},
		{"bad link", func(p *handler.RequestAuthEmailRegister) { p.Link = "invalid" }},
		{"short password", func(p *handler.RequestAuthEmailRegister) { p.Password = "123" }},
		{"unknown account reset", func(p *handler.RequestAuthEmailRegister) { p.Name = "" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := handler.RequestAuthEmailRegister{Name: "Reader", Email: challenge.Email, Password: "valid-password", Code: challenge.Code}
			tc.edit(&p)
			flowRequest(t, router, "POST", "/auth/email/register", p, "", 400)
			var saved entity.UserEmailVerify
			require.NoError(t, app.Dao().DB().First(&saved, challenge.ID).Error)
			require.Equal(t, challenge.Code, saved.Code)
			require.Zero(t, saved.Attempts)
		})
	}
	flowRequest(t, router, "POST", "/auth/email/login", handler.RequestAuthEmailLogin{Email: challenge.Email}, "", 400)
	flowRequest(t, router, "POST", "/auth/email/login", handler.RequestAuthEmailLogin{Email: challenge.Email, Password: "wrong"}, "", 401)
	flowRequest(t, router, "POST", "/auth/email/login", handler.RequestAuthEmailLogin{Email: challenge.Email, Code: challenge.Code}, "", 401)
	app.Conf().Auth.Email.Enabled = false
	flowRequest(t, router, "POST", "/auth/email/register", map[string]any{}, "", 400)
	flowRequest(t, router, "POST", "/auth/email/login", map[string]any{}, "", 400)
}
