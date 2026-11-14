package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAdminUserManagement(t *testing.T) {
	app, router := NewApiTestApp()
	cached, err := cache.New(config.CacheConf{Type: config.CacheTypeBuiltin})
	require.NoError(t, err)
	t.Cleanup(cached.Close)
	app.Dao().SetCache(dao.NewCacheAdaptor(cached))
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.UserCreate(app.App, router)
	handler.UserUpdate(app.App, router)
	handler.UserDelete(app.App, router)
	handler.UserList(app.App, router)

	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)

	create := handler.ParamsUserCreate{
		Name: "Managed reader", Email: "managed@example.org", Password: "initial-password",
		ReceiveEmail: true,
	}
	flowRequest(t, router, http.MethodPost, "/users", create, "", http.StatusForbidden)
	var created handler.ResponseUserCreate
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/users", create, admin, http.StatusOK), &created))
	flowRequest(t, router, http.MethodPost, "/users", create, admin, http.StatusBadRequest)
	require.NotZero(t, created.ID)
	require.False(t, created.IsAdmin)
	stored := app.Dao().FindUserByID(created.ID)
	require.True(t, stored.CheckPassword(create.Password))

	path := fmt.Sprintf("/users/%d", created.ID)
	update := handler.ParamsUserUpdate{
		Name: create.Name, Email: create.Email, Password: "new-password",
		BadgeName: "Moderator", IsAdmin: true, ReceiveEmail: false,
	}
	flowRequest(t, router, http.MethodPut, path, update, admin, http.StatusOK)
	stored = app.Dao().FindUserByID(created.ID)
	require.False(t, stored.CheckPassword(create.Password))
	require.True(t, stored.CheckPassword(update.Password))
	require.False(t, stored.ReceiveEmail)
	require.Equal(t, update.BadgeName, stored.BadgeName)

	managed, err := common.LoginGetUserToken(stored, app.Conf().AppKey, 60)
	require.NoError(t, err)
	flowRequest(t, router, http.MethodGet, "/users/admin", nil, managed, http.StatusOK)
	flowRequest(t, router, http.MethodDelete, path, nil, admin, http.StatusOK)
	flowRequest(t, router, http.MethodGet, "/users/admin", nil, managed, http.StatusForbidden)
	require.True(t, app.Dao().FindUserByID(created.ID).IsEmpty())
}

func TestAdminUserValidationAndStorageFailures(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.UserCreate(app.App, router)
	handler.UserUpdate(app.App, router)
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	original := app.Dao().FindUserByID(1001)
	other := app.Dao().FindUserByID(1002)
	for _, invalid := range []struct{ name, email, link string }{
		{other.Name, other.Email, ""},
		{"Reader", "invalid-email", ""},
		{"Reader", "reader@example.org", "invalid-url"},
	} {
		create := handler.ParamsUserCreate{Name: invalid.name, Email: invalid.email, Link: invalid.link}
		flowRequest(t, router, "POST", "/users", create, admin, 400)
		update := handler.ParamsUserUpdate{Name: invalid.name, Email: invalid.email, Link: invalid.link}
		flowRequest(t, router, "PUT", "/users/1001", update, admin, 400)
	}
	update := handler.ParamsUserUpdate{Name: "Reader", Email: "reader@example.org"}
	flowRequest(t, router, "PUT", "/users/999999", update, admin, 404)
	const callback = "test:user_write_failure"
	fail := func(tx *gorm.DB) { tx.AddError(errors.New("user storage unavailable")) }
	require.NoError(t, app.Dao().DB().Callback().Create().Before("gorm:create").Register(callback, fail))
	t.Cleanup(func() { require.NoError(t, app.Dao().DB().Callback().Create().Remove(callback)) })
	flowRequest(t, router, "POST", "/users", handler.ParamsUserCreate{Name: update.Name, Email: update.Email}, admin, 500)
	require.True(t, app.Dao().FindUser(update.Name, update.Email).IsEmpty())
	require.NoError(t, app.Dao().DB().Callback().Update().Before("gorm:update").Register(callback, fail))
	t.Cleanup(func() { require.NoError(t, app.Dao().DB().Callback().Update().Remove(callback)) })
	flowRequest(t, router, "PUT", "/users/1001", update, admin, 500)
	require.Equal(t, original, app.Dao().FindUserByID(original.ID))
}
