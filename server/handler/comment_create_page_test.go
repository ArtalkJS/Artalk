package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
)

func TestRejectedCommentDoesNotCreatePage(t *testing.T) {
	app, router := NewApiTestApp()
	defer app.Cleanup()
	app.Conf().Captcha.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.CommentCreate(app.App, router)
	for _, tc := range []struct {
		name              string
		rid               uint
		token             string
		anonymousDisabled bool
		status            int
	}{
		{"invalid token", 0, "invalid", false, 401},
		{"missing parent", 999999, "", false, 404},
		{"wrong page", 1000, "", false, 400},
		{"anonymous disabled", 0, "", true, 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			key := "/rejected/" + tc.name
			app.Conf().Auth.Enabled = tc.anonymousDisabled
			app.Conf().Auth.Anonymous = !tc.anonymousDisabled
			body, err := json.Marshal(handler.ParamsCommentCreate{
				Name: "Visitor", Email: "visitor@example.com", Content: "Test", SiteName: "Site A", PageKey: key, Rid: tc.rid,
			})
			require.NoError(t, err)
			req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			resp, err := router.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, tc.status, resp.StatusCode)
			require.True(t, app.Dao().FindPage(key, "Site A").IsEmpty())
		})
	}
}
