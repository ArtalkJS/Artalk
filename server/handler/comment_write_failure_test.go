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

func TestApprovalRetriesAfterNotificationWriteFailure(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing=%v", existing), func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			handler.CommentUpdate(app.App, router)
			d := app.Dao()
			parent := entity.Comment{Content: "Parent", UserID: 1001, SiteName: "Site A", PageKey: "/test/1000.html"}
			require.NoError(t, d.CreateComment(&parent))
			reply := entity.Comment{Content: "Pending reply", UserID: 1002, SiteName: parent.SiteName, PageKey: parent.PageKey, Rid: parent.ID, RootID: parent.ID, IsPending: true}
			require.NoError(t, d.CreateComment(&reply))
			if existing {
				n := d.NewNotify(parent.UserID, reply.ID)
				require.NoError(t, d.NotifySetRead(&n))
			}
			reject := func(tx *gorm.DB) {
				if tx.Statement.Table == "atk_notifies" {
					tx.AddError(errors.New("notification write failed"))
				}
			}
			require.NoError(t, d.DB().Callback().Create().Before("gorm:create").Register("test:notify_failure", reject))
			require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:notify_failure", reject))
			t.Cleanup(func() {
				d.DB().Callback().Create().Remove("test:notify_failure")
				d.DB().Callback().Update().Remove("test:notify_failure")
			})
			admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
			require.NoError(t, err)
			path := fmt.Sprintf("/comments/%d", reply.ID)
			update := handler.ParamsCommentUpdate{Content: reply.Content, PageKey: reply.PageKey, SiteName: reply.SiteName, Rid: reply.Rid}
			flowRequest(t, router, http.MethodPut, path, update, admin, http.StatusInternalServerError)
			require.True(t, d.FindComment(reply.ID).IsPending)
			notify := d.FindNotify(parent.UserID, reply.ID)
			if existing {
				require.True(t, notify.IsRead)
			} else {
				require.True(t, notify.IsEmpty())
			}
			require.NoError(t, d.DB().Callback().Create().Remove("test:notify_failure"))
			require.NoError(t, d.DB().Callback().Update().Remove("test:notify_failure"))
			flowRequest(t, router, http.MethodPut, path, update, admin, http.StatusOK)
			require.False(t, d.FindComment(reply.ID).IsPending)
			notify = d.FindNotify(parent.UserID, reply.ID)
			require.False(t, notify.IsEmpty())
			require.False(t, notify.IsRead)
		})
	}
}

func TestCommentEditChecksRelatedWrites(t *testing.T) {
	for _, tc := range []struct {
		name, table string
		edit        func(*handler.ParamsCommentUpdate)
	}{
		{"new author", "atk_users", func(p *handler.ParamsCommentUpdate) { p.Nick = "New author"; p.Email = "new@example.org" }},
		{"author link", "atk_users", func(p *handler.ParamsCommentUpdate) { p.Link = "https://changed.example.org" }},
		{"new page", "atk_pages", func(p *handler.ParamsCommentUpdate) { p.PageKey = "/new-page" }},
		{"successful edit", "", func(p *handler.ParamsCommentUpdate) {
			p.Nick = "New author"
			p.Email = "new@example.org"
			p.Link = "https://changed.example.org"
			p.PageKey = "/new-page"
			p.IP = "192.0.2.1"
			p.UA = "Test browser"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			handler.CommentUpdate(app.App, router)
			d := app.Dao()
			comment := entity.Comment{Content: "Original", UserID: 1001, PageKey: "/test/1000.html", SiteName: "Site A"}
			require.NoError(t, d.CreateComment(&comment))
			admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
			require.NoError(t, err)
			reject := func(tx *gorm.DB) {
				if tx.Statement.Table == tc.table {
					tx.AddError(errors.New("related write failed"))
				}
			}
			require.NoError(t, d.DB().Callback().Create().Before("gorm:create").Register("test:related_write", reject))
			require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:related_write", reject))
			t.Cleanup(func() {
				d.DB().Callback().Create().Remove("test:related_write")
				d.DB().Callback().Update().Remove("test:related_write")
			})
			p := handler.ParamsCommentUpdate{Content: "Edited", PageKey: comment.PageKey, SiteName: comment.SiteName}
			tc.edit(&p)
			status := http.StatusInternalServerError
			if tc.table == "" {
				status = http.StatusOK
			}
			flowRequest(t, router, http.MethodPut, fmt.Sprintf("/comments/%d", comment.ID), p, admin, status)
			saved := d.FindComment(comment.ID)
			if tc.table != "" {
				require.Equal(t, comment.Content, saved.Content)
				require.Equal(t, comment.PageKey, saved.PageKey)
				require.Equal(t, comment.UserID, saved.UserID)
			} else {
				require.Equal(t, p.Content, saved.Content)
				require.Equal(t, p.PageKey, saved.PageKey)
				require.Equal(t, p.UA, saved.UA)
				require.Equal(t, p.IP, saved.IP)
				require.Equal(t, p.Email, d.FindUserByID(saved.UserID).Email)
			}
		})
	}
}

func TestCommentEditKeepsLegacyAuthor(t *testing.T) {
	for _, email := range []string{"", "invalid-email"} {
		t.Run(email, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			handler.CommentUpdate(app.App, router)
			d := app.Dao()
			author, err := d.NewUser("Legacy reader", email, "")
			require.NoError(t, err)
			comment := entity.Comment{Content: "Old comment", UserID: author.ID, SiteName: "Site A", PageKey: "/legacy", IsPending: true}
			require.NoError(t, d.CreateComment(&comment))
			admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
			require.NoError(t, err)
			flowRequest(t, router, http.MethodPut, fmt.Sprintf("/comments/%d", comment.ID), handler.ParamsCommentUpdate{
				Content: "Edited", PageKey: comment.PageKey, SiteName: comment.SiteName,
			}, admin, http.StatusOK)
			saved := d.FindComment(comment.ID)
			require.Equal(t, author.ID, saved.UserID)
			require.Equal(t, "Edited", saved.Content)
			require.False(t, saved.IsPending)
		})
	}
}
