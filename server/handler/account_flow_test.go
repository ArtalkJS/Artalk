package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

func TestProfileEmailChangeAndStorageFailure(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.UserInfoUpdate(app.App, router)
	user := app.Dao().FindUserByID(1001)
	token, err := common.LoginGetUserToken(user, app.Conf().AppKey, 60)
	require.NoError(t, err)
	profile := handler.RequestUserInfoUpdate{Name: user.Name, Email: "not-an-email"}
	flowRequest(t, router, "POST", "/user", profile, token, 400)
	profile.Email = app.Dao().FindUserByID(1002).Email
	flowRequest(t, router, "POST", "/user", profile, token, 400)
	profile.Name = app.Dao().FindUserByID(1002).Name
	flowRequest(t, router, "POST", "/user", profile, token, 400)
	require.Equal(t, user, app.Dao().FindUserByID(user.ID))

	challenge := entity.UserEmailVerify{Email: "new-address@example.org", Code: "654321", ExpiresAt: time.Now().Add(time.Minute)}
	require.NoError(t, app.Dao().DB().Create(&challenge).Error)
	profile = handler.RequestUserInfoUpdate{Name: user.Name, Email: challenge.Email, Code: challenge.Code}
	flowRequest(t, router, "POST", "/user", profile, token, 200)
	require.Equal(t, challenge.Email, app.Dao().FindUserByID(user.ID).Email)
	profile.Link = "reader.example.org"
	flowRequest(t, router, "POST", "/user", profile, token, 200)
	require.Equal(t, "https://reader.example.org", app.Dao().FindUserByID(user.ID).Link)
	const callback = "test:profile_save_failure"
	require.NoError(t, app.Dao().DB().Callback().Update().Before("gorm:update").Register(callback, func(tx *gorm.DB) {
		tx.AddError(errors.New("profile storage unavailable"))
	}))
	t.Cleanup(func() { require.NoError(t, app.Dao().DB().Callback().Update().Remove(callback)) })
	profile.Name = "Unsaved name"
	flowRequest(t, router, "POST", "/user", profile, token, 500)
	require.Equal(t, user.Name, app.Dao().FindUserByID(user.ID).Name)
}
