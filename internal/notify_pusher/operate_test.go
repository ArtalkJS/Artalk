package notify_pusher_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestAdminPushHonorsModerationSettings(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	d := app.Dao()

	type message struct {
		Subject string            `json:"notify_subject"`
		Comment entity.CommentDTO `json:"comment"`
	}
	requests := make(chan message, 4)
	callback := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requests <- msg
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(callback.Close)
	settings := config.AdminNotifyConf{
		Email:         &config.AdminEmailConf{},
		NotifySubject: "New comment",
		WebHook:       config.NotifyWebHookConf{Enabled: true, URL: callback.URL},
	}
	comment := entity.Comment{Content: "Review me", UserID: 1001, SiteName: "Site A", PageKey: "/test/1000.html", IsPending: true}
	require.NoError(t, d.CreateComment(&comment))

	pusher := notify_pusher.NewNotifyPusher(&notify_pusher.NotifyPusherConf{Dao: d, AdminNotifyConf: settings})
	pusher.Push(&comment, nil)
	require.Empty(t, requests)

	settings.NotifyPending = true
	pusher = notify_pusher.NewNotifyPusher(&notify_pusher.NotifyPusherConf{Dao: d, AdminNotifyConf: settings})
	pusher.Push(&comment, nil)
	require.Len(t, requests, 1)
	request := <-requests
	require.Equal(t, settings.NotifySubject, request.Subject)
	require.Equal(t, comment.ID, request.Comment.ID)
	require.Equal(t, comment.Content, request.Comment.Content)
	require.True(t, request.Comment.IsPending)
}

func TestAdminEmailRecipientRules(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	d := app.Dao()
	admin := d.FindUserByID(1000)
	require.NoError(t, d.DB().Model(&entity.User{}).Where("id <> ?", admin.ID).Update("is_admin", false).Error)
	duplicate := entity.User{Name: "Other admin profile", Email: admin.Email, IsAdmin: true, ReceiveEmail: true}
	require.NoError(t, d.CreateUser(&duplicate))
	for _, tc := range []struct {
		name                                            string
		author, parent                                  uint
		pending, notifyPending, noise, disabled, optOut bool
		want                                            []uint
	}{
		{name: "root comment", author: 1001, want: []uint{admin.ID}},
		{name: "pending", author: 1001, pending: true},
		{name: "notify pending", author: 1001, pending: true, notifyPending: true, want: []uint{admin.ID}},
		{name: "administrator comment", author: admin.ID},
		{name: "administrator replies", author: admin.ID, parent: admin.ID, noise: true},
		{name: "reply to administrator", author: 1001, parent: admin.ID, noise: true, want: []uint{admin.ID}},
		{name: "reader conversation", author: 1001, parent: 1002, want: []uint{1002}},
		{name: "reader conversation with noise mode", author: 1001, parent: 1002, noise: true, want: []uint{1002, admin.ID}},
		{name: "administrator opted out", author: 1001, optOut: true},
		{name: "administrator email disabled", author: 1001, parent: admin.ID, disabled: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			admin.ReceiveEmail = !tc.optOut
			require.NoError(t, d.UpdateUser(&admin))
			var delivered []uint
			pusher := notify_pusher.NewNotifyPusher(&notify_pusher.NotifyPusherConf{
				Dao: d, AdminNotifyConf: config.AdminNotifyConf{Email: &config.AdminEmailConf{Enabled: !tc.disabled}, NoiseMode: tc.noise, NotifyPending: tc.notifyPending},
				EmailPush: func(n *entity.Notify) error { delivered = append(delivered, n.UserID); return nil },
			})
			var parent *entity.Comment
			if tc.parent != 0 {
				parent = &entity.Comment{UserID: tc.parent, Content: "Parent"}
				require.NoError(t, d.CreateComment(parent))
			}
			comment := entity.Comment{UserID: tc.author, Content: "Reply", IsPending: tc.pending}
			if parent != nil {
				comment.Rid = parent.ID
			}
			require.NoError(t, d.CreateComment(&comment))
			pusher.Push(&comment, parent)
			require.ElementsMatch(t, tc.want, delivered)
		})
	}
}
