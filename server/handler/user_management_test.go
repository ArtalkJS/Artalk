package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestAdminUserManagement(t *testing.T) {
	app, router := NewApiTestApp()
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
