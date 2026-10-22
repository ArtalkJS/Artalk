package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestPageListHasComments(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.PageList(app.App, router)
	for i, key := range []string{"public", "pending", "empty", "deleted", "cross-site"} {
		page := app.Dao().NewPage("/filter/"+key, "Filter "+key, "Site A")
		// Deterministic ordering for pagination assertions.
		require.NoError(t, app.Dao().DB().Model(&page).Update("created_at", time.Date(2026, 1, i+1, 0, 0, 0, 0, time.UTC)).Error)
		if key == "empty" {
			continue
		}
		comment := entity.Comment{PageKey: page.Key, SiteName: page.SiteName, UserID: 1001, Content: "test", IsPending: key == "pending"}
		if key == "cross-site" {
			comment.SiteName = "Site B"
		}
		require.NoError(t, app.Dao().DB().Create(&comment).Error)
		if key == "deleted" {
			require.NoError(t, app.Dao().DB().Delete(&comment).Error)
		}
	}
	token, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	for _, tc := range []struct {
		query string
		total int64
		keys  []string
	}{
		{"&has_comments=true", 2, []string{"/filter/pending", "/filter/public"}},
		{"&has_comments=true&limit=1", 2, []string{"/filter/pending"}},
		{"&has_comments=true&limit=1&offset=1", 2, []string{"/filter/public"}},
		{"&has_comments=true&offset=2", 2, []string{}},
		{"", 5, []string{"/filter/cross-site", "/filter/deleted", "/filter/empty", "/filter/pending", "/filter/public"}},
		{"&has_comments=false", 5, []string{"/filter/cross-site", "/filter/deleted", "/filter/empty", "/filter/pending", "/filter/public"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/pages?site_name=Site%20A&search=Filter"+tc.query, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := router.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)
			var result handler.ResponsePageList
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
			require.Equal(t, tc.total, result.Total)
			keys := make([]string, 0, len(result.Pages))
			for _, page := range result.Pages {
				keys = append(keys, page.Key)
			}
			require.Equal(t, tc.keys, keys)
			require.NotNil(t, result.Pages)
		})
	}
	// Filtering must not expose the admin endpoint to anonymous callers.
	resp, err := router.Test(httptest.NewRequest("GET", "/pages?has_comments=true", nil))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 403, resp.StatusCode)
}
