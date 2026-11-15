package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCommentCreationWriteFailures(t *testing.T) {
	for _, table := range []string{"atk_users", "atk_pages", "atk_comments"} {
		t.Run(table, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			app.Conf().Captcha.Enabled = false
			router.Use(limiter.ActionLimitMiddleware(app.App, limiter.ActionLimitConf{}))
			handler.CommentCreate(app.App, router)
			storage := app.Dao().DB()
			var before int64
			require.NoError(t, storage.Model(&entity.Comment{}).Count(&before).Error)
			require.NoError(t, storage.Callback().Create().Before("gorm:create").Register("test:write_failure", func(tx *gorm.DB) {
				if tx.Statement.Table == table {
					tx.AddError(errors.New("simulated storage failure"))
				}
			}))
			t.Cleanup(func() { require.NoError(t, storage.Callback().Create().Remove("test:write_failure")) })
			flowRequest(t, router, http.MethodPost, "/comments", handler.ParamsCommentCreate{
				Name: "New visitor", Email: "new@example.org", Content: "Must not appear", PageKey: "/write-failure", SiteName: "Site A",
			}, "", http.StatusInternalServerError)
			var after int64
			require.NoError(t, storage.Model(&entity.Comment{}).Count(&after).Error)
			require.Equal(t, before, after)
		})
	}
}

func TestFailedApprovalDoesNotPublishOrNotify(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.CommentUpdate(app.App, router)
	d := app.Dao()
	parent := entity.Comment{Content: "Parent", UserID: 1001, SiteName: "Site A", PageKey: "/test/1000.html"}
	require.NoError(t, d.CreateComment(&parent))
	reply := entity.Comment{Content: "Pending reply", UserID: 1002, SiteName: parent.SiteName, PageKey: parent.PageKey, Rid: parent.ID, RootID: parent.ID, IsPending: true}
	require.NoError(t, d.CreateComment(&reply))
	admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:approval_failure", func(tx *gorm.DB) {
		if tx.Statement.Table == "atk_comments" {
			tx.AddError(errors.New("cannot save approval"))
		}
	}))
	t.Cleanup(func() { require.NoError(t, d.DB().Callback().Update().Remove("test:approval_failure")) })
	flowRequest(t, router, http.MethodPut, fmt.Sprintf("/comments/%d", reply.ID), handler.ParamsCommentUpdate{
		Content: "Approved content", PageKey: reply.PageKey, SiteName: reply.SiteName, Rid: reply.Rid,
	}, admin, http.StatusInternalServerError)
	saved := d.FindComment(reply.ID)
	require.True(t, saved.IsPending)
	require.Equal(t, reply.Content, saved.Content)
	require.True(t, d.FindNotify(parent.UserID, reply.ID).IsEmpty(), "failed approval must not notify the reader")
}
