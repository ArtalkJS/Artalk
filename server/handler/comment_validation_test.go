package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
)

func TestCommentValidationDoesNotWrite(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().Captcha.Enabled = false
	router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
	handler.CommentCreate(app.App, router)
	handler.CommentUpdate(app.App, router)
	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	parent := entity.Comment{Content: "Parent", UserID: 1001, PageKey: "/validation", SiteName: "Site B"}
	require.NoError(t, app.Dao().CreateComment(&parent))
	pending := entity.Comment{Content: "Pending", UserID: 1001, PageKey: parent.PageKey, SiteName: "Site A", IsPending: true}
	require.NoError(t, app.Dao().CreateComment(&pending))
	for _, tc := range []struct {
		name   string
		edit   func(*handler.ParamsCommentCreate)
		status int
	}{
		{"invalid email", func(p *handler.ParamsCommentCreate) { p.Email = "invalid" }, 400},
		{"invalid link", func(p *handler.ParamsCommentCreate) { p.Link = "javascript:alert(1)" }, 400},
		{"missing content", func(p *handler.ParamsCommentCreate) { p.Content = "" }, 400},
		{"missing site", func(p *handler.ParamsCommentCreate) { p.SiteName = "Unknown" }, 404},
		{"pending parent", func(p *handler.ParamsCommentCreate) { p.Rid = pending.ID }, 400},
		{"different site", func(p *handler.ParamsCommentCreate) { p.Rid = parent.ID }, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := handler.ParamsCommentCreate{Name: "Reader", Email: "reader@example.org", Content: "Reply", PageKey: parent.PageKey, SiteName: "Site A"}
			tc.edit(&p)
			flowRequest(t, router, http.MethodPost, "/comments", p, "", tc.status)
		})
	}
	for _, tc := range []struct {
		name   string
		edit   func(*handler.ParamsCommentUpdate)
		status int
	}{
		{"invalid email", func(p *handler.ParamsCommentUpdate) { p.Email = "invalid" }, 400},
		{"invalid link", func(p *handler.ParamsCommentUpdate) { p.Link = "invalid" }, 400},
		{"missing site", func(p *handler.ParamsCommentUpdate) { p.SiteName = "Unknown" }, 404},
		{"missing content", func(p *handler.ParamsCommentUpdate) { p.Content = "" }, 400},
	} {
		t.Run("update "+tc.name, func(t *testing.T) {
			p := handler.ParamsCommentUpdate{Content: pending.Content, PageKey: pending.PageKey, SiteName: pending.SiteName, IsPending: true}
			tc.edit(&p)
			flowRequest(t, router, http.MethodPut, fmt.Sprintf("/comments/%d", pending.ID), p, admin, tc.status)
		})
	}
	require.Equal(t, pending.Content, app.Dao().FindComment(pending.ID).Content)
	require.True(t, app.Dao().FindComment(pending.ID).IsPending)
}
