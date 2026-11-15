package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestMergeMovesAccountDataToSelectedProfile(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	app.Conf().LoginTimeout = 3600
	handler.AuthMergeCheck(app.App, router)
	handler.AuthMergeApply(app.App, router)
	handler.UserInfo(app.App, router)

	d := app.Dao()
	const email = "merge@example.org"
	source, err := d.NewUser("Old profile", email, "")
	require.NoError(t, err)
	target, err := d.NewUser("Preferred profile", email, "")
	require.NoError(t, err)
	comment := entity.Comment{UserID: source.ID, Content: "Comment to keep", SiteName: "Site A", PageKey: "/test/1000.html"}
	require.NoError(t, d.CreateComment(&comment))
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
	require.Equal(t, notify.ID, d.FindNotify(target.ID, comment.ID).ID)
	require.NoError(t, d.DB().First(&vote, vote.ID).Error)
	require.Equal(t, target.ID, vote.UserID)
	require.Equal(t, target.ID, d.FindAuthIdentityByRemoteUID(identity.Provider, identity.RemoteUID).UserID)

	flowRequest(t, router, http.MethodGet, "/auth/merge", nil, oldToken, http.StatusUnauthorized)
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/auth/merge", nil, result.UserToken, http.StatusOK), &check))
	require.False(t, check.NeedMerge)
}
