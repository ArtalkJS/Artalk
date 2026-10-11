package handler_test

import (
	"net/http"
	"net/http/httptest"
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
