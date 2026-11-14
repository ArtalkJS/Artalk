package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestDeleteCommentRemovesThreadAndAssociatedData(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.CommentDelete(app.App, router)
	d := app.Dao()

	makeComment := func(parentID uint) entity.Comment {
		t.Helper()
		c := entity.Comment{
			UserID: 1001, Rid: parentID, Content: "A thread comment",
			SiteName: "Site A", PageKey: "/test/1000.html",
		}
		require.NoError(t, d.CreateComment(&c))
		return c
	}
	root := makeComment(0)
	reply := makeComment(root.ID)
	nested := makeComment(reply.ID)
	other := makeComment(0)
	for _, id := range []uint{root.ID, reply.ID, nested.ID, other.ID} {
		require.NotZero(t, d.NewNotify(1002, id).ID)
		require.NoError(t, d.DB().Create(&entity.Vote{TargetID: id, Type: entity.VoteTypeCommentUp, UserID: 1002}).Error)
	}

	admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	path := fmt.Sprintf("/comments/%d", root.ID)
	flowRequest(t, router, http.MethodDelete, path, nil, "", http.StatusForbidden)
	require.False(t, d.FindComment(root.ID).IsEmpty())
	flowRequest(t, router, http.MethodDelete, path, nil, admin, http.StatusOK)
	flowRequest(t, router, http.MethodDelete, path, nil, admin, http.StatusNotFound)
	for _, id := range []uint{root.ID, reply.ID, nested.ID} {
		require.True(t, d.FindComment(id).IsEmpty())
		require.True(t, d.FindNotify(1002, id).IsEmpty())
		require.Zero(t, d.GetVoteNum(id, string(entity.VoteTypeCommentUp)))
	}
	require.False(t, d.FindComment(other.ID).IsEmpty())
	require.False(t, d.FindNotify(1002, other.ID).IsEmpty())
	require.EqualValues(t, 1, d.GetVoteNum(other.ID, string(entity.VoteTypeCommentUp)))
}
