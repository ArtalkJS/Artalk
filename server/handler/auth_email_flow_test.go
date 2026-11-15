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
