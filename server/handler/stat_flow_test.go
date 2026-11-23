package handler_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestMostCommentedPagesCountOnlyVisibleCommentsFromTheSameSite(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.Stat(app.App, router)

	dao := app.Dao()
	for _, page := range []entity.Page{
		{SiteName: "blog", Key: "/shared"},
		{SiteName: "blog", Key: "/local"},
		{SiteName: "wiki", Key: "/shared"},
	} {
		require.NoError(t, dao.CreatePage(&page))
	}
	for _, comment := range []entity.Comment{
		{SiteName: "blog", PageKey: "/local"},
		{SiteName: "blog", PageKey: "/shared", IsPending: true},
		{SiteName: "blog", PageKey: "/shared", IsPending: true},
		{SiteName: "wiki", PageKey: "/shared"},
		{SiteName: "wiki", PageKey: "/shared"},
	} {
		require.NoError(t, dao.CreateComment(&comment))
	}

	var result struct {
		Data []entity.PageDTO `json:"data"`
	}
	body := flowRequest(t, router, http.MethodGet, "/stats/comment_most_pages?site_name=blog&limit=1", nil, "", http.StatusOK)
	require.NoError(t, json.Unmarshal(body, &result))
	require.Len(t, result.Data, 1)
	require.Equal(t, "blog", result.Data[0].SiteName)
	require.Equal(t, "/local", result.Data[0].Key)
}

func TestPageStatisticsReflectWritesWithinEachSite(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.PagePV(app.App, router)
	handler.Stat(app.App, router)

	const key = "/shared-statistics-page"
	logView := func(site string) {
		t.Helper()
		flowRequest(t, router, http.MethodPost, "/pages/pv",
			handler.ParamsPagePV{SiteName: site, PageKey: key}, "", http.StatusOK)
	}
	logView("Site A")
	logView("Site A")
	logView("Site B")

	for _, c := range []entity.Comment{
		{SiteName: "Site A", PageKey: key, UserID: 1001, Content: "Visible A"},
		{SiteName: "Site A", PageKey: key, UserID: 1002, Content: "Pending A", IsPending: true},
		{SiteName: "Site B", PageKey: key, UserID: 1001, Content: "Visible B"},
		{SiteName: "Site B", PageKey: key, UserID: 1002, Content: "Another visible B"},
	} {
		require.NoError(t, app.Dao().CreateComment(&c))
	}

	for _, tc := range []struct {
		site        string
		views       int
		numComments int
	}{
		{"Site A", 2, 1},
		{"Site B", 1, 2},
	} {
		t.Run(tc.site, func(t *testing.T) {
			const missing = "/missing-statistics-page"
			query := url.Values{"site_name": {tc.site}, "page_keys": {key + "," + missing}}
			for endpoint, want := range map[string]int{"page_pv": tc.views, "page_comment": tc.numComments} {
				path := "/stats/" + endpoint + "?" + query.Encode()
				var body struct {
					Data map[string]int `json:"data"`
				}
				require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, path, nil, "", http.StatusOK), &body))
				require.Equal(t, map[string]int{key: want, missing: 0}, body.Data, endpoint)
			}
		})
	}
}
