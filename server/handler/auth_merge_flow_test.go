package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMergeMovesAccountDataToSelectedProfile(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().LoginTimeout = 3600
	handler.AuthMergeCheck(app.App, router)
	handler.AuthMergeApply(app.App, router)
	handler.UserInfo(app.App, router)

	d := app.Dao()
	cached, err := cache.New(config.CacheConf{Type: config.CacheTypeBuiltin})
	require.NoError(t, err)
	t.Cleanup(cached.Close)
	d.SetCache(dao.NewCacheAdaptor(cached))
	const email = "merge@example.org"
	source, err := d.NewUser("Old profile", email, "")
	require.NoError(t, err)
	target, err := d.NewUser("Preferred profile", email, "")
	require.NoError(t, err)
	comment := entity.Comment{UserID: source.ID, Content: "Comment to keep", SiteName: "Site A", PageKey: "/test/1000.html"}
	require.NoError(t, d.CreateComment(&comment))
	unrelated := entity.Comment{UserID: 1001, Content: "Unrelated cache entry"}
	require.NoError(t, d.CreateComment(&unrelated))
	require.NoError(t, d.DB().Model(&unrelated).UpdateColumn("content", "Changed directly").Error)
	require.Equal(t, "Unrelated cache entry", d.FindComment(unrelated.ID).Content)
	notify := d.NewNotify(source.ID, comment.ID)
	require.NotZero(t, notify.ID)
	vote := entity.Vote{UserID: source.ID, TargetID: comment.ID, Type: entity.VoteTypeCommentUp}
	require.NoError(t, d.DB().Create(&vote).Error)
	identity := entity.AuthIdentity{Provider: "fixture", RemoteUID: "merge-reader", UserID: source.ID}
	require.NoError(t, d.DB().Create(&identity).Error)

	oldToken, err := common.LoginGetUserToken(source, app.Conf().AppKey, 60)
	require.NoError(t, err)
	var check handler.ResponseAuthDataMergeCheck
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/auth/merge", nil, oldToken, http.StatusOK), &check))
	require.True(t, check.NeedMerge)
	require.ElementsMatch(t, []string{source.Name, target.Name}, check.UserNames)

	request := handler.RequestAuthDataMergeApply{UserName: target.Name}
	flowRequest(t, router, http.MethodPost, "/auth/merge", request, "", http.StatusUnauthorized)
	var result handler.ResponseAuthDataMergeApply
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodPost, "/auth/merge", request, oldToken, http.StatusOK), &result))
	require.EqualValues(t, 1, result.UpdatedComment)
	require.EqualValues(t, 1, result.UpdatedNotify)
	require.EqualValues(t, 1, result.UpdatedVote)
	require.EqualValues(t, 1, result.DeletedUser)
	require.NotEmpty(t, result.UserToken)

	require.Equal(t, target.ID, d.FindComment(comment.ID).UserID)
	require.Equal(t, "Unrelated cache entry", d.FindComment(unrelated.ID).Content)
	require.Equal(t, notify.ID, d.FindNotify(target.ID, comment.ID).ID)
	require.NoError(t, d.DB().First(&vote, vote.ID).Error)
	require.Equal(t, target.ID, vote.UserID)
	require.Equal(t, target.ID, d.FindAuthIdentityByRemoteUID(identity.Provider, identity.RemoteUID).UserID)

	flowRequest(t, router, http.MethodGet, "/auth/merge", nil, oldToken, http.StatusUnauthorized)
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/auth/merge", nil, result.UserToken, http.StatusOK), &check))
	require.False(t, check.NeedMerge)
}

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

func TestAccountMergeRollsBackFailedDataTransfers(t *testing.T) {
	for _, table := range []string{"atk_comments", "atk_notifies", "atk_votes", "atk_auth_identities"} {
		t.Run(table, func(t *testing.T) {
			app, router := NewApiTestApp()
			t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
			handler.AuthMergeApply(app.App, router)
			d := app.Dao()
			source, err := d.NewUser("Source", "rollback@example.org", "")
			require.NoError(t, err)
			target, err := d.NewUser("Target", source.Email, "")
			require.NoError(t, err)
			comment := entity.Comment{UserID: source.ID, Content: "Keep ownership"}
			require.NoError(t, d.CreateComment(&comment))
			notify := d.NewNotify(source.ID, comment.ID)
			vote, err := d.NewVote(comment.ID, entity.VoteTypeCommentUp, source.ID, "", "")
			require.NoError(t, err)
			identity := entity.AuthIdentity{UserID: source.ID, Provider: "fixture", RemoteUID: "rollback"}
			require.NoError(t, d.CreateAuthIdentity(&identity))
			require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:merge_failure", func(tx *gorm.DB) {
				if tx.Statement.Table == table {
					tx.AddError(errors.New("merge storage failure"))
				}
			}))
			t.Cleanup(func() { require.NoError(t, d.DB().Callback().Update().Remove("test:merge_failure")) })
			token, err := common.LoginGetUserToken(source, app.Conf().AppKey, 60)
			require.NoError(t, err)
			flowRequest(t, router, "POST", "/auth/merge", handler.RequestAuthDataMergeApply{UserName: target.Name}, token, 500)
			for _, record := range []any{&comment, &notify, &vote, &identity} {
				require.NoError(t, d.DB().First(record).Error)
			}
			require.Equal(t, source.ID, comment.UserID)
			require.Equal(t, source.ID, notify.UserID)
			require.Equal(t, source.ID, vote.UserID)
			require.Equal(t, source.ID, identity.UserID)
			require.False(t, d.FindUserByID(source.ID).IsEmpty())
		})
	}
}
