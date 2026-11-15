package dao_test

import (
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMigrateRootIDsForExistingThreads(t *testing.T) {
	storage, err := db.NewDB(config.DBConf{Type: config.TypeSQLite, Dsn: filepath.Join(t.TempDir(), "legacy.db"), TablePrefix: "atk_"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })

	type oldComment struct {
		gorm.Model
		Rid     uint
		Content string
	}
	const table = "atk_comments"
	require.NoError(t, storage.Table(table).AutoMigrate(&oldComment{}))
	records := []oldComment{
		{Model: gorm.Model{ID: 1}, Content: "Thread A"},
		{Model: gorm.Model{ID: 2}, Rid: 1, Content: "Reply"},
		{Model: gorm.Model{ID: 3}, Rid: 2, Content: "Nested reply"},
		{Model: gorm.Model{ID: 4}, Content: "Thread B"},
		{Model: gorm.Model{ID: 5}, Rid: 4, Content: "Other reply"},
	}
	require.NoError(t, storage.Table(table).Create(&records).Error)
	require.False(t, storage.Migrator().HasColumn(table, "root_id"))

	d := dao.NewDao(storage)
	for _, pass := range []string{"first migration", "rerun"} {
		t.Run(pass, func(t *testing.T) {
			var comments []entity.Comment
			require.NoError(t, storage.Order("id").Find(&comments).Error)
			require.Len(t, comments, len(records))
			wantRoots := []uint{0, 1, 1, 0, 4}
			for i, comment := range comments {
				require.Equal(t, records[i].ID, comment.ID)
				require.Equal(t, records[i].Content, comment.Content)
				require.Equal(t, records[i].Rid, comment.Rid)
				require.Equal(t, wantRoots[i], comment.RootID)
			}
		})
		d.MigrateModels()
	}
}
