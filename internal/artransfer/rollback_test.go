package artransfer

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImportRollsBackRelatedWrites(t *testing.T) {
	for _, table := range []string{"sites", "pages", "users", "comments", "votes"} {
		t.Run(table, func(t *testing.T) {
			storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "import.db")})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
			d := dao.NewDao(storage)
			require.NoError(t, storage.Callback().Create().Before("gorm:create").Register("test:import_failure", func(tx *gorm.DB) {
				if tx.Statement.Table == table {
					tx.AddError(errors.New("import storage failure"))
				}
			}))
			data, err := json.Marshal([]entity.Artran{{ID: "1", Content: "Import", Nick: "Reader", Email: "reader@example.org", SiteName: "New site", PageKey: "/new", VoteUp: "1"}})
			require.NoError(t, err)
			err = RunImportArtrans(d, &ImportParams{JsonData: string(data), Assumeyes: true}, func(string) {})
			require.Error(t, err)
			for _, model := range []any{&entity.Site{}, &entity.Page{}, &entity.User{}, &entity.Comment{}, &entity.Vote{}} {
				var count int64
				require.NoError(t, storage.Model(model).Count(&count).Error)
				require.Zero(t, count, "partial import persisted %T", model)
			}
		})
	}
}

func TestSkippedImportRowsPreserveReplyLinksAndDates(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "skip.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
	d := dao.NewDao(storage)
	rows := []entity.Artran{
		{ID: "missing-site", PageKey: "/post", Content: "Skip"},
		{ID: "missing-page", SiteName: "Blog", Content: "Skip"},
		{ID: "parent", SiteName: "Blog", PageKey: "/post", Content: "Parent", Email: "invalid", Link: "reader.example.org", CreatedAt: "2024-03-04T05:06:07Z"},
		{ID: "reply", Rid: "parent", SiteName: "Blog", PageKey: "/post", Content: "Reply", Nick: "Reader", Email: "reader@example.org", CreatedAt: "2024-03-05T05:06:07Z"},
	}
	data, err := json.Marshal(rows)
	require.NoError(t, err)
	require.NoError(t, RunImportArtrans(d, &ImportParams{JsonData: string(data), Assumeyes: true}, func(string) {}))
	var saved []entity.Comment
	require.NoError(t, storage.Order("id").Find(&saved).Error)
	require.Len(t, saved, 2)
	require.Zero(t, saved[0].Rid)
	require.Equal(t, saved[0].ID, saved[1].Rid)
	require.Equal(t, saved[0].ID, saved[1].RootID)
	for i, comment := range saved {
		require.Equal(t, rows[i+2].CreatedAt, comment.CreatedAt.UTC().Format(time.RFC3339))
	}
	user := d.FindUserByID(saved[0].UserID)
	require.Equal(t, "Anonymous", user.Name)
	require.Equal(t, "anonymous@example.org", user.Email)
	require.Equal(t, "https://reader.example.org", user.Link)
}

func TestSkippedAncestorsKeepImportedRepliesVisible(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "ancestors.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })
	d := dao.NewDao(storage)
	rows := []entity.Artran{
		{ID: "missing-root", PageKey: "/post"},
		{ID: "parent", Rid: "missing-root", SiteName: "Blog", PageKey: "/post", Content: "Parent"},
		{ID: "child", Rid: "parent", SiteName: "Blog", PageKey: "/post", Content: "Child"},
		{ID: "root", SiteName: "Blog", PageKey: "/post", Content: "Root"},
		{ID: "missing-middle", Rid: "root", SiteName: "Blog"},
		{ID: "other-parent", Rid: "missing-middle", SiteName: "Blog", PageKey: "/post", Content: "Other parent"},
		{ID: "other-child", Rid: "other-parent", SiteName: "Blog", PageKey: "/post", Content: "Other child"},
	}
	data, err := json.Marshal(rows)
	require.NoError(t, err)
	require.NoError(t, RunImportArtrans(d, &ImportParams{JsonData: string(data), Assumeyes: true}, func(string) {}))
	var saved []entity.Comment
	require.NoError(t, storage.Order("id").Find(&saved).Error)
	require.Len(t, saved, 5)
	for _, pair := range [][2]int{{0, 1}, {3, 4}} {
		parent, child := saved[pair[0]], saved[pair[1]]
		require.Zero(t, parent.Rid)
		require.Zero(t, parent.RootID)
		require.Equal(t, parent.ID, child.Rid)
		require.Equal(t, parent.ID, child.RootID)
	}
	require.Zero(t, saved[2].Rid)
}
