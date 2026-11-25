package dao_test

import (
	"fmt"
	"os"
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

// Use SQLite by default, or a disposable external database for compatibility checks.
func newMigrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	conf := config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "migration.db"), TablePrefix: "atk_"}
	if kind := os.Getenv("ARTALK_TEST_DB_TYPE"); kind != "" {
		conf.Type = config.DBType(kind)
		conf.Dsn = os.Getenv("ARTALK_TEST_DB_DSN")
		require.NotEmpty(t, conf.Dsn)
		conf.TablePrefix = fmt.Sprintf("migration_%x_", time.Now().UnixNano())
	}
	storage, err := db.NewDB(conf)
	require.NoError(t, err)
	t.Cleanup(func() {
		if conf.Type != config.TypeSQLite {
			for _, model := range []any{&entity.Notify{}, &entity.Vote{}, &entity.Comment{}, &entity.UserEmailVerify{}, &entity.AuthIdentity{}, &entity.Page{}, &entity.Site{}, &entity.User{}} {
				require.NoError(t, storage.Migrator().DropTable(model))
			}
		}
		require.NoError(t, db.CloseDB(storage))
	})
	return storage
}

func TestMigrationAddsAnnotationColumns(t *testing.T) {
	storage := newMigrationDB(t)
	require.NoError(t, storage.Exec(`CREATE TABLE atk_comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		page_key TEXT,
		site_name TEXT,
		user_id INTEGER,
		rid INTEGER,
		root_id INTEGER,
		is_collapsed INTEGER DEFAULT 0,
		is_pending INTEGER DEFAULT 0,
		is_pinned INTEGER DEFAULT 0,
		vote_up INTEGER DEFAULT 0,
		vote_down INTEGER DEFAULT 0,
		is_verified INTEGER DEFAULT 0,
		ua TEXT,
		ip TEXT,
		guest_token_hash TEXT DEFAULT '',
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	)`).Error)
	require.NoError(t, storage.Exec(`INSERT INTO atk_comments (content, page_key, site_name, user_id) VALUES ('legacy', '/legacy', 'Site A', 1001)`).Error)

	d := dao.NewDao(storage)
	require.True(t, storage.Migrator().HasColumn(&entity.Comment{}, "anchor_exact"))
	require.True(t, storage.Migrator().HasColumn(&entity.Comment{}, "is_annotation_thread"))

	var comment entity.Comment
	require.NoError(t, d.DB().Where("content = ?", "legacy").First(&comment).Error)
	require.False(t, comment.IsAnnotationThread)
	require.False(t, comment.AnchorResolved)
	require.Equal(t, "", comment.AnchorExact)
	require.Equal(t, -1, comment.AnchorStart)
	require.Equal(t, -1, comment.AnchorEnd)
}
