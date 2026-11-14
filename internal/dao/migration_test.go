package dao_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
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
