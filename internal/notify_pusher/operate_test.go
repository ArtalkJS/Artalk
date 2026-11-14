package notify_pusher_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/notify_pusher"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
)

func TestReplyNotificationsRespectRecipientPreferences(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	d := app.Dao()

	const readerID, otherReaderID, adminID = 1001, 1002, 1000
	var delivered []entity.Notify
	pusher := notify_pusher.NewNotifyPusher(&notify_pusher.NotifyPusherConf{
		AdminNotifyConf: config.AdminNotifyConf{Email: &config.AdminEmailConf{}},
		Dao:             d,
		EmailPush: func(notify *entity.Notify) error {
			delivered = append(delivered, *notify)
			return nil
		},
	})

	for _, tc := range []struct {
		name        string
		replyAuthor uint
		emailOn     bool
		pending     bool
		wantEmail   bool
	}{
		{name: "approved reply", replyAuthor: adminID, emailOn: true, wantEmail: true},
		{name: "pending reply", replyAuthor: otherReaderID, emailOn: true, pending: true},
		{name: "self reply", replyAuthor: readerID, emailOn: true},
		{name: "recipient opted out", replyAuthor: otherReaderID},
	} {
		t.Run(tc.name, func(t *testing.T) {
			delivered = nil
			reader := d.FindUserByID(readerID)
			reader.ReceiveEmail = tc.emailOn
			require.NoError(t, d.UpdateUser(&reader))
			parent := entity.Comment{UserID: readerID, Content: "Parent", SiteName: "Site A", PageKey: "/test/1000.html"}
			require.NoError(t, d.CreateComment(&parent))
			reply := entity.Comment{
				UserID: tc.replyAuthor, Content: "Reply", Rid: parent.ID, RootID: parent.ID,
				SiteName: parent.SiteName, PageKey: parent.PageKey, IsPending: tc.pending,
			}
			require.NoError(t, d.CreateComment(&reply))
			pusher.Push(&reply, &parent)

			saved := d.FindNotify(readerID, reply.ID)
			if !tc.wantEmail {
				require.Empty(t, delivered)
				require.True(t, saved.IsEmpty())
				return
			}
			require.Len(t, delivered, 1)
			require.NotZero(t, saved.ID)
			require.Equal(t, saved.ID, delivered[0].ID)
			require.EqualValues(t, readerID, delivered[0].UserID)
			require.False(t, saved.IsRead)
		})
	}
}
