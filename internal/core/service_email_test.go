package core

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
)

type emailTransport func(*http.Request) (*http.Response, error)

func (f emailTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestEmailServiceDelivery(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "email.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
	app := NewApp(&config.Config{SiteDefault: "Artalk"})
	app.SetDao(dao.NewDao(storage))
	cfg := &app.Conf().Email
	cfg.Enabled = true
	cfg.SendType = config.TypeAliDM
	cfg.SendAddr = "noreply@example.org"
	cfg.SendName = "Artalk"
	cfg.MailSubject = "Reader reply"
	cfg.MailTpl = filepath.Join(t.TempDir(), "reader.html")
	require.NoError(t, os.WriteFile(cfg.MailTpl, []byte("Reader template"), 0600))
	app.Conf().AdminNotify.Email = &config.AdminEmailConf{MailSubject: "Admin reply"}
	app.Conf().AdminNotify.Email.MailTpl = filepath.Join(t.TempDir(), "admin.html")
	require.NoError(t, os.WriteFile(app.Conf().AdminNotify.Email.MailTpl, []byte("Admin template"), 0600))
	received := make(chan url.Values, 3)
	previous := http.DefaultTransport
	http.DefaultTransport = emailTransport(func(r *http.Request) (*http.Response, error) {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		received <- r.Form
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"EnvId":"fixture"}`))}, nil
	})
	t.Cleanup(func() { http.DefaultTransport = previous })
	service := NewEmailService(app)
	require.NoError(t, service.Init())
	t.Cleanup(func() { require.NoError(t, service.Dispose()) })
	for _, admin := range []bool{false, true} {
		user := entity.User{Name: "Reader", Email: "reader@example.org", IsAdmin: admin}
		require.NoError(t, app.Dao().CreateUser(&user))
		comment := entity.Comment{Content: "A reply", UserID: user.ID}
		require.NoError(t, app.Dao().CreateComment(&comment))
		notification := app.Dao().NewNotify(user.ID, comment.ID)
		service.AsyncSend(&notification)
		select {
		case sent := <-received:
			require.Equal(t, user.Email, sent.Get("ToAddress"))
			require.Equal(t, map[bool]string{false: "Reader reply", true: "Admin reply"}[admin], sent.Get("Subject"))
			require.Equal(t, map[bool]string{false: "Reader template", true: "Admin template"}[admin], sent.Get("HtmlBody"))
		case <-time.After(5 * time.Second):
			t.Fatal("email was not delivered")
		}
		require.Eventually(t, func() bool {
			var saved entity.Notify
			return storage.First(&saved, notification.ID).Error == nil && saved.IsEmailed && saved.EmailAt != nil
		}, 5*time.Second, time.Millisecond)
	}
	service.AsyncSendTo("Login code", "123456", "login@example.org")
	select {
	case sent := <-received:
		require.Equal(t, "Login code", sent.Get("Subject"))
		require.Equal(t, "login@example.org", sent.Get("ToAddress"))
	case <-time.After(5 * time.Second):
		t.Fatal("login email was not delivered")
	}
	cfg.Enabled = false
	service.AsyncSendTo("Disabled", "unused", "login@example.org")
	service.AsyncSend(nil)
	require.Empty(t, received)
}
