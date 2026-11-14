package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestPageRenameMovesOnlyCommentsInItsSite(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.PageUpdate(app.App, router)

	const oldKey = "/same-page-key"
	page := app.Dao().FindCreatePage(oldKey, "", "Site A")
	other := app.Dao().FindCreatePage(oldKey, "", "Site B")
	require.NotZero(t, page.ID)
	require.NotZero(t, other.ID)
	first := entity.Comment{Content: "Site A", PageKey: oldKey, SiteName: page.SiteName, UserID: 1001}
	second := entity.Comment{Content: "Site B", PageKey: oldKey, SiteName: other.SiteName, UserID: 1001}
	require.NoError(t, app.Dao().DB().Create(&first).Error)
	require.NoError(t, app.Dao().DB().Create(&second).Error)
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)

	update := handler.ParamsPageUpdate{SiteName: page.SiteName, Key: "/moved", Title: "Moved"}
	flowRequest(t, router, http.MethodPut, fmt.Sprintf("/pages/%d", page.ID), update, admin, http.StatusOK)
	require.Equal(t, page.ID, app.Dao().FindPage(update.Key, page.SiteName).ID)
	require.Equal(t, update.Key, app.Dao().FindComment(first.ID).PageKey)
	require.Equal(t, oldKey, app.Dao().FindComment(second.ID).PageKey)
	require.Equal(t, other.ID, app.Dao().FindPage(oldKey, other.SiteName).ID)
}
