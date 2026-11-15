package handler_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAccountMergeWithoutCacheAvoidsFullCommentScan(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.AuthMergeApply(app.App, router)
	d := app.Dao()
	d.SetCache(nil)
	source, err := d.NewUser("Source", "scan@example.org", "")
	require.NoError(t, err)
	target, err := d.NewUser("Target", source.Email, "")
	require.NoError(t, err)
	fullScan := false
	require.NoError(t, d.DB().Callback().Query().After("gorm:query").Register("test:full_scan", func(tx *gorm.DB) {
		if tx.Statement.Table == "atk_comments" && !strings.Contains(strings.ToLower(tx.Statement.SQL.String()), "user_id") {
			fullScan = true
		}
	}))
	t.Cleanup(func() { d.DB().Callback().Query().Remove("test:full_scan") })
	token, err := common.LoginGetUserToken(source, app.Conf().AppKey, 60)
	require.NoError(t, err)
	flowRequest(t, router, http.MethodPost, "/auth/merge", handler.RequestAuthDataMergeApply{UserName: target.Name}, token, http.StatusOK)
	require.False(t, fullScan)
}
