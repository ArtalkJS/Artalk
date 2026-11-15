package dao_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMigrateRootIDsForExistingThreads(t *testing.T) {
	for _, fallback := range []bool{false, true} {
		t.Run(map[bool]string{false: "recursive SQL", true: "legacy fallback"}[fallback], func(t *testing.T) {
			testMigrateRootIDs(t, fallback)
		})
	}
}

func testMigrateRootIDs(t *testing.T, fallback bool) {
	t.Helper()
	storage := newMigrationDB(t)

	type oldComment struct {
		gorm.Model
		Rid     uint
		Content string
	}
	table := storage.NamingStrategy.TableName("Comment")
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

	if fallback {
		require.NoError(t, storage.Callback().Row().Before("gorm:row").Register("test:no_cte", func(tx *gorm.DB) {
			if strings.HasPrefix(tx.Statement.SQL.String(), "WITH RECURSIVE") {
				tx.AddError(errors.New("recursive CTE unavailable"))
			}
		}))
	}
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
