package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
)

func TestApprovalDeliversPendingReplyNotification(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.CommentUpdate(app.App, router)
	handler.NotifyList(app.App, router)
	handler.NotifyReadAll(app.App, router)

	parent := entity.Comment{Content: "Parent", UserID: 1001, SiteName: "Site A", PageKey: "/test/1000.html"}
	require.NoError(t, app.Dao().DB().Create(&parent).Error)
	reply := entity.Comment{
		Content: "Held for review", UserID: 1002, SiteName: parent.SiteName,
		PageKey: parent.PageKey, Rid: parent.ID, RootID: parent.ID, IsPending: true,
	}
	require.NoError(t, app.Dao().DB().Create(&reply).Error)
	notify, err := core.AppService[*core.NotifyService](app.App)
	require.NoError(t, err)
	require.NoError(t, notify.Push(&reply, &parent))
	require.True(t, app.Dao().FindNotify(parent.UserID, reply.ID).IsEmpty())

	admin, err := common.LoginGetUserToken(app.Dao().FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	reader, err := common.LoginGetUserToken(app.Dao().FindUserByID(parent.UserID), app.Conf().AppKey, 60)
	require.NoError(t, err)
	update := handler.ParamsCommentUpdate{SiteName: reply.SiteName, PageKey: reply.PageKey, Content: reply.Content, Rid: parent.ID}
	path := fmt.Sprintf("/comments/%d", reply.ID)
	flowRequest(t, router, http.MethodPut, path, update, reader, http.StatusForbidden)
	require.True(t, app.Dao().FindComment(reply.ID).IsPending)
	flowRequest(t, router, http.MethodPut, path, update, admin, http.StatusOK)
	require.False(t, app.Dao().FindComment(reply.ID).IsPending)

	var unread handler.ResponseNotifyList
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/notifies", nil, reader, http.StatusOK), &unread))
	require.Contains(t, notificationCommentIDs(unread), reply.ID)

	flowRequest(t, router, http.MethodPost, "/notifies/read", nil, reader, http.StatusOK)
	require.NoError(t, json.Unmarshal(flowRequest(t, router, http.MethodGet, "/notifies", nil, reader, http.StatusOK), &unread))
	require.Zero(t, unread.Count)
}

func notificationCommentIDs(list handler.ResponseNotifyList) []uint {
	ids := make([]uint, 0, len(list.Notifies))
	for _, item := range list.Notifies {
		ids = append(ids, item.CommentID)
	}
	return ids
}
