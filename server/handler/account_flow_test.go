package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
)

func TestPasswordLoginAndProfile(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Captcha.Enabled = false
	app.Conf().LoginTimeout = 3600
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.UserLogin(app.App, router)
	handler.UserInfoUpdate(app.App, router)
	handler.UserList(app.App, router)

	const password = "a-test-password"
	reader := entity.User{Name: "Reader", Email: "reader@example.org"}
	require.NoError(t, reader.SetPasswordEncrypt(password))
	require.NoError(t, app.Dao().DB().Create(&reader).Error)
	other := entity.User{Name: "Other reader", Email: reader.Email}
	require.NoError(t, app.Dao().DB().Create(&other).Error)

	login := handler.ParamsUserLogin{Email: reader.Email, Password: password}
	var selection struct {
		Names []string `json:"need_name_select"`
	}
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/user/access_token", login, "", http.StatusBadRequest), &selection))
	require.ElementsMatch(t, []string{reader.Name, other.Name}, selection.Names)

	login.Name = reader.Name
	login.Password = "wrong-password"
	flowRequest(t, router, http.MethodPost, "/user/access_token", login, "", http.StatusUnauthorized)
	login.Password = password
	var signedIn handler.ResponseUserLogin
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/user/access_token", login, "", http.StatusOK), &signedIn))
	require.NotEmpty(t, signedIn.Token)
	require.Equal(t, reader.ID, signedIn.User.ID)
	flowRequest(t, router, http.MethodGet, "/users/all", nil, signedIn.Token, http.StatusForbidden)

	profile := handler.RequestUserInfoUpdate{Email: reader.Email, Name: "Updated reader", Link: "https://reader.example.org"}
	flowRequest(t, router, http.MethodPost, "/user", profile, "", http.StatusUnauthorized)
	var updated handler.ResponseUserInfoUpdate
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/user", profile, signedIn.Token, http.StatusOK), &updated))
	require.Equal(t, reader.ID, updated.User.ID)
	require.Equal(t, profile.Name, updated.User.Name)
	require.Equal(t, profile.Link, app.Dao().FindUserByID(reader.ID).Link)
	require.Equal(t, other.Name, app.Dao().FindUserByID(other.ID).Name)

	// Users with passwords stored by older releases can still sign in.
	legacy := app.Dao().FindUserByID(1000)
	var admin handler.ResponseUserLogin
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/user/access_token",
		handler.ParamsUserLogin{Email: legacy.Email, Password: "123456"}, "", http.StatusOK), &admin))
	var admins handler.ResponseAdminUserList
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/users/admin", nil, admin.Token, http.StatusOK), &admins))
	require.NotEmpty(t, admins.Users)
	for _, user := range admins.Users {
		require.True(t, user.IsAdmin)
		require.NotEqual(t, reader.ID, user.ID)
	}
}
