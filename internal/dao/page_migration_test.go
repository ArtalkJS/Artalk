package dao_test

import (
	"errors"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMergeLegacyPages(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "preserve metadata", true: "rollback on write failure"}[fail], func(t *testing.T) {
			storage := newMigrationDB(t)
			d := dao.NewDao(storage)
			require.NoError(t, storage.Exec("ALTER TABLE "+d.GetTableName(&entity.Page{})+" ADD COLUMN accessible_url TEXT").Error)
			pages := []entity.Page{
				{SiteName: "ab", Key: "c", Title: "Original", AdminOnly: true, PV: 3, VoteUp: 2},
				{SiteName: "ab", Key: "c", Title: "Duplicate", PV: 7, VoteUp: 1, VoteDown: 4},
				{SiteName: "a", Key: "bc", Title: "Different site and key", PV: 11},
			}
			require.NoError(t, storage.Create(&pages).Error)
			if fail {
				require.NoError(t, storage.Callback().Create().Before("gorm:create").Register("test:merge_failure", func(tx *gorm.DB) {
					tx.AddError(errors.New("cannot write merged pages"))
				}))
				require.ErrorContains(t, d.MergePages(), "cannot write merged pages")
				var saved []entity.Page
				require.NoError(t, storage.Order("id").Find(&saved).Error)
				require.Equal(t, pages, saved)
				require.True(t, storage.Migrator().HasColumn(&entity.Page{}, "accessible_url"))
				return
			}
			for range 2 {
				require.NoError(t, d.MergePages())
				var saved []entity.Page
				require.NoError(t, storage.Order("id").Find(&saved).Error)
				require.Len(t, saved, 2)
				require.Equal(t, pages[0].ID, saved[0].ID)
				require.Equal(t, "Original", saved[0].Title)
				require.True(t, saved[0].AdminOnly)
				require.Equal(t, 10, saved[0].PV)
				require.Equal(t, 3, saved[0].VoteUp)
				require.Equal(t, 4, saved[0].VoteDown)
				require.Equal(t, pages[2], saved[1])
				require.False(t, storage.Migrator().HasColumn(&entity.Page{}, "accessible_url"))
			}
		})
	}
}
