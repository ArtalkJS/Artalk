package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentListRejectsInvalidScope(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	var caughtPanic any
	// The test catches panics so a regression cannot terminate the test process.
	router.Use(func(c *fiber.Ctx) (err error) {
		defer func() {
			if p := recover(); p != nil {
				caughtPanic = p
				err = c.SendStatus(http.StatusInternalServerError)
			}
		}()
		return c.Next()
	})
	handler.CommentList(app.App, router)

	resp, err := router.Test(httptest.NewRequest(http.MethodGet, "/comments?page_key=/test&scope=invalid", nil))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Nil(t, caughtPanic, "an untrusted query must not panic")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUnavailableIPDatabaseDoesNotHideComments(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.CommentList(app.App, router)
	const url = "/comments?page_key=/test/1000.html&site_name=Site%20A&flat_mode=true&limit=100"
	var before, after handler.ResponseCommentList
	require.NoError(t, json.Unmarshal(flowRequest(t, router, "GET", url, nil, "", 200), &before))
	require.NotEmpty(t, before.Comments)
	app.Conf().IPRegion.Enabled = true
	app.Conf().IPRegion.DBPath = filepath.Join(t.TempDir(), "missing.ipdb")
	require.NoError(t, json.Unmarshal(flowRequest(t, router, "GET", url, nil, "", 200), &after))
	require.Equal(t, before, after)
}
