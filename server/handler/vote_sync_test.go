package handler_test

import (
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestVoteSyncRepairsCountersFromStoredVotes(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.VoteSync(app.App, router)
	d := app.Dao()

	page := d.FindCreatePage("/vote-sync", "Vote sync", "Site A")
	require.NotZero(t, page.ID)
	page.VoteUp, page.VoteDown = 5, 9
	require.NoError(t, d.UpdatePage(&page))
	comment := entity.Comment{UserID: 1001, Content: "Vote target", SiteName: page.SiteName, PageKey: page.Key, VoteUp: 7, VoteDown: 9}
	require.NoError(t, d.CreateComment(&comment))
	for _, vote := range []entity.Vote{
		{TargetID: page.ID, Type: entity.VoteTypePageDown, UserID: 1001},
		{TargetID: comment.ID, Type: entity.VoteTypeCommentUp, UserID: 1001},
		{TargetID: comment.ID, Type: entity.VoteTypeCommentUp, UserID: 1002},
	} {
		require.NoError(t, d.DB().Create(&vote).Error)
	}

	admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	flowRequest(t, router, http.MethodPost, "/votes/sync", nil, "", http.StatusForbidden)
	require.EqualValues(t, 5, d.FindPageByID(page.ID).VoteUp)
	flowRequest(t, router, http.MethodPost, "/votes/sync", nil, admin, http.StatusOK)

	page = d.FindPageByID(page.ID)
	comment = d.FindComment(comment.ID)
	require.Zero(t, page.VoteUp)
	require.EqualValues(t, 1, page.VoteDown)
	require.Equal(t, 2, comment.VoteUp)
	require.Zero(t, comment.VoteDown)
}
