package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentGetVisibility(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	handler.CommentGet(app.App, router)
	parent := app.Dao().FindComment(1001)
	parent.IsPending = true
	require.NoError(t, app.Dao().DB().Save(&parent).Error)

	for _, tc := range []struct {
		name          string
		id            uint
		userID        uint
		status        int
		parentVisible bool
	}{
		{"anonymous pending", 1001, 0, 404, false},
		{"unrelated user pending", 1001, 1002, 404, false},
		{"owner pending", 1001, 1001, 200, true},
		{"admin pending", 1001, 1000, 200, true},
		{"anonymous reply to pending", 1002, 0, 200, false},
		{"owner reply to pending", 1002, 1001, 200, true},
		{"admin reply to pending", 1002, 1000, 200, true},
		{"anonymous public", 1000, 0, 200, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/comments/%d", tc.id), nil)
			if tc.userID != 0 {
				token, err := common.LoginGetUserToken(app.Dao().FindUserByID(tc.userID), app.Conf().AppKey, 60)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
			}
			resp, err := router.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tc.status, resp.StatusCode)
			if tc.status == 200 {
				var result handler.ResponseCommentGet
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
				assert.Equal(t, tc.parentVisible, result.ReplyComment != nil)
			}
		})
	}

	// Deleted parents must not appear as empty comment objects.
	require.NoError(t, app.Dao().DB().Delete(&entity.Comment{}, 1001).Error)
	resp, err := router.Test(httptest.NewRequest(http.MethodGet, "/comments/1002", nil))
	require.NoError(t, err)
	defer resp.Body.Close()
	var result handler.ResponseCommentGet
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Nil(t, result.ReplyComment)
}
