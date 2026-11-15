package artransfer

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

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
