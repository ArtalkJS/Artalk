package core_test

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
)

type moderationTransport func(*http.Request) (*http.Response, error)

func (f moderationTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestModerationServicePersistsProviderDecisions(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "moderation.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
	d := dao.NewDao(storage)
	app := core.NewApp(&config.Config{})
	app.SetDao(d)
	cfg := &app.Conf().Moderator
	cfg.AkismetKey = "fixture-key"
	reader, err := d.NewUser("Reader", "reader@example.org", "")
	require.NoError(t, err)

	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	for _, tc := range []struct {
		name, response string
		failBlock      bool
		pending        bool
	}{
		{"clean comment", "false", false, false},
		{"spam comment", "true", false, true},
		{"provider error with fail open", "service error", false, false},
		{"provider error with fail closed", "service error", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg.ApiFailBlock = tc.failBlock
			service := core.NewAntiSpamService(app)
			require.NoError(t, service.Init())
			t.Cleanup(func() { require.NoError(t, service.Dispose()) })
			comment := entity.Comment{Content: "A comment to moderate", PageKey: "/moderated", UserID: reader.ID}
			require.NoError(t, d.CreateComment(&comment))
			calls := 0
			http.DefaultTransport = moderationTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "https://fixture-key.rest.akismet.com/1.1/comment-check", r.URL.String())
				if err := r.ParseForm(); err != nil {
					return nil, fmt.Errorf("parse moderation request: %w", err)
				}
				require.Equal(t, "https://blog.example.org", r.Form.Get("blog"))
				require.Equal(t, reader.Email, r.Form.Get("comment_author_email"))
				require.Equal(t, comment.Content, r.Form.Get("comment_content"))
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(tc.response))}, nil
			})

			service.CheckAndBlock(&core.AntiSpamCheckPayload{Comment: &comment, ReqReferer: "https://blog.example.org/post"})
			require.Equal(t, 1, calls)
			saved := d.FindComment(comment.ID)
			require.Equal(t, tc.pending, saved.IsPending)
			require.Equal(t, comment.Content, saved.Content)
		})
	}
}
