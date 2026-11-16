package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func flowRequest(t *testing.T, router *fiber.App, method, path string, body any, token string, status int) []byte {
	t.Helper()
	res, err := securityRequest(router, method, path, body, token, "")
	require.NoError(t, err)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, status, res.StatusCode, "%s %s: %s", method, path, data)
	return data
}

func TestSiteManagementKeepsPagesAndCommentsInTheirSite(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	c, err := cache.New(config.CacheConf{Type: config.CacheTypeBuiltin})
	require.NoError(t, err)
	t.Cleanup(c.Close)
	app.Dao().SetCache(dao.NewCacheAdaptor(c))

	handler.SiteCreate(app.App, router)
	handler.SiteUpdate(app.App, router)
	handler.SiteList(app.App, router)
	handler.SiteDelete(app.App, router)
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	reader, err := common.LoginGetUserToken(app.Dao().FindUserByID(1001), app.Conf().AppKey, 60)
	require.NoError(t, err)

	create := handler.ParamsSiteCreate{Name: "Flow site", Urls: []string{"https://flow.example.org"}}
	flowRequest(t, router, http.MethodPost, "/sites", create, reader, http.StatusForbidden)
	require.True(t, app.Dao().FindSite(create.Name).IsEmpty())
	var site entity.SiteDTO
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/sites", create, admin, http.StatusOK), &site))
	require.Equal(t, create.Urls, site.Urls)
	flowRequest(t, router, http.MethodPost, "/sites", create, admin, http.StatusBadRequest)
	flowRequest(t, router, http.MethodPut, fmt.Sprintf("/sites/%d", site.ID), map[string]any{"name": site.Name, "urls": nil}, admin, http.StatusOK)
	require.Empty(t, app.Dao().FindSiteByID(site.ID).Urls)

	const pageKey = "/shared"
	page := app.Dao().FindCreatePage(pageKey, "", site.Name)
	other := app.Dao().FindCreatePage(pageKey, "", "Site A")
	require.NotZero(t, page.ID)
	require.NotZero(t, other.ID)
	comment := entity.Comment{Content: "Move with the page", PageKey: pageKey, SiteName: site.Name, UserID: 1001}
	require.NoError(t, app.Dao().DB().Create(&comment).Error)
	require.Equal(t, page.ID, app.Dao().FindPage(pageKey, site.Name).ID)
	require.Equal(t, site.Name, app.Dao().FindComment(comment.ID).SiteName)

	path := fmt.Sprintf("/sites/%d", site.ID)
	update := handler.ParamsSiteUpdate{Name: "Site A", Urls: create.Urls}
	flowRequest(t, router, http.MethodPut, path, update, admin, http.StatusBadRequest)
	update.Name = "Renamed flow site"
	flowRequest(t, router, http.MethodPut, path, update, admin, http.StatusOK)

	var list handler.ResponseSiteList
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/sites", nil, admin, http.StatusOK), &list))
	require.Equal(t, len(list.Sites), list.Count)
	var names []string
	for _, listed := range list.Sites {
		names = append(names, listed.Name)
	}
	require.Contains(t, names, update.Name)
	require.NotContains(t, names, site.Name)
	require.True(t, app.Dao().FindSite(site.Name).IsEmpty())
	require.True(t, app.Dao().FindPage(pageKey, site.Name).IsEmpty())
	require.Equal(t, page.ID, app.Dao().FindPage(pageKey, update.Name).ID)
	require.Equal(t, update.Name, app.Dao().FindComment(comment.ID).SiteName)

	flowRequest(t, router, http.MethodDelete, path, nil, reader, http.StatusForbidden)
	flowRequest(t, router, http.MethodDelete, path, nil, admin, http.StatusOK)
	require.True(t, app.Dao().FindSite(update.Name).IsEmpty())
	require.True(t, app.Dao().FindPage(pageKey, update.Name).IsEmpty())
	require.True(t, app.Dao().FindComment(comment.ID).IsEmpty())
	require.Equal(t, other.ID, app.Dao().FindPage(pageKey, other.SiteName).ID)
}
