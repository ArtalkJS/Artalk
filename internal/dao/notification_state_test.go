package dao_test

import (
	"errors"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNotificationReadStateRemainsUserScoped(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	d := app.Dao()
	parent := entity.Comment{Content: "Parent", UserID: 1001}
	require.NoError(t, d.CreateComment(&parent))
	reply := entity.Comment{Content: "Reply", UserID: 1002, Rid: parent.ID}
	require.NoError(t, d.CreateComment(&reply))
	own := d.NewNotify(parent.UserID, reply.ID)
	other := d.NewNotify(reply.UserID, parent.ID)
	require.Equal(t, parent.ID, d.FindNotifyParentComment(&own).ID)
	require.True(t, d.FindNotifyParentComment(&other).IsEmpty())
	require.Empty(t, d.FindUnreadNotifies(0))
	require.Error(t, d.UserNotifyMarkAllAsRead(0))
	require.NotEmpty(t, d.FindUnreadNotifies(parent.UserID))
	require.NoError(t, d.UserNotifyMarkAllAsRead(parent.UserID))
	require.Empty(t, d.FindUnreadNotifies(parent.UserID))
	require.False(t, d.FindNotify(reply.UserID, parent.ID).IsRead)
	require.NoError(t, d.NotifySetInitial(&own))
	require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:read_failure", func(tx *gorm.DB) { tx.AddError(errors.New("cannot persist read state")) }))
	t.Cleanup(func() { require.NoError(t, d.DB().Callback().Update().Remove("test:read_failure")) })
	require.ErrorContains(t, d.UserNotifyMarkAllAsRead(parent.UserID), "cannot persist read state")
	require.False(t, d.FindNotify(parent.UserID, reply.ID).IsRead)
}
