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
	for _, failure := range []string{"none", "pages", "votes", "delete"} {
		t.Run(failure, func(t *testing.T) {
			storage := newMigrationDB(t)
			d := dao.NewDao(storage)
			require.NoError(t, storage.Table(d.GetTableName(&entity.Page{})).Migrator().AddColumn(&struct {
				AccessibleURL string
			}{}, "accessible_url"))
			pages := []entity.Page{
				{SiteName: "ab", Key: "c", Title: "Original", AdminOnly: true, PV: 3, VoteUp: 2},
				{SiteName: "ab", Key: "c", Title: "Duplicate", PV: 7, VoteUp: 1, VoteDown: 4},
				{SiteName: "a", Key: "bc", Title: "Different site and key", PV: 11},
			}
			require.NoError(t, storage.Create(&pages).Error)
			deleted := entity.Page{SiteName: "Blog", Key: "/deleted", Title: "Keep deletion history"}
			require.NoError(t, storage.Create(&deleted).Error)
			require.NoError(t, storage.Delete(&deleted).Error)
			votes := []entity.Vote{
				{TargetID: pages[1].ID, Type: entity.VoteTypePageUp},
				{TargetID: pages[1].ID, Type: entity.VoteTypePageDown},
				{TargetID: pages[1].ID, Type: entity.VoteTypeCommentUp},
			}
			require.NoError(t, storage.Create(&votes).Error)
			// Compare database snapshots to account for driver timestamp precision and time zones.
			// Select mapped fields because the migration removes a legacy column.
			pageQuery := storage.Session(&gorm.Session{QueryFields: true})
			require.NoError(t, pageQuery.Order("id").Find(&pages).Error)
			require.NoError(t, pageQuery.Unscoped().First(&deleted, deleted.ID).Error)
			require.NoError(t, storage.Order("id").Find(&votes).Error)
			if failure != "none" {
				reject := func(tx *gorm.DB) {
					tx.AddError(errors.New("cannot merge pages"))
				}
				if failure == "delete" {
					require.NoError(t, storage.Callback().Delete().Before("gorm:delete").Register("test:merge_failure", reject))
				} else {
					require.NoError(t, storage.Callback().Update().Before("gorm:update").Register("test:merge_failure", func(tx *gorm.DB) {
						model := any(&entity.Page{})
						if failure == "votes" {
							model = &entity.Vote{}
						}
						if tx.Statement.Table == d.GetTableName(model) {
							reject(tx)
						}
					}))
				}
				require.ErrorContains(t, d.MergePages(), "cannot merge pages")
				var saved []entity.Page
				require.NoError(t, pageQuery.Order("id").Find(&saved).Error)
				require.Equal(t, pages, saved)
				var savedVotes []entity.Vote
				require.NoError(t, storage.Order("id").Find(&savedVotes).Error)
				require.Equal(t, votes, savedVotes)
				require.True(t, storage.Migrator().HasColumn(&entity.Page{}, "accessible_url"))
				return
			}
			for range 2 {
				require.NoError(t, d.MergePages())
				var saved []entity.Page
				require.NoError(t, pageQuery.Order("id").Find(&saved).Error)
				require.Len(t, saved, 2)
				require.Equal(t, pages[0].ID, saved[0].ID)
				require.Equal(t, "Original", saved[0].Title)
				require.True(t, saved[0].AdminOnly)
				require.Equal(t, 10, saved[0].PV)
				require.Equal(t, 3, saved[0].VoteUp)
				require.Equal(t, 4, saved[0].VoteDown)
				require.Equal(t, pages[2], saved[1])
				require.False(t, storage.Migrator().HasColumn(&entity.Page{}, "accessible_url"))
				var kept entity.Page
				require.NoError(t, pageQuery.Unscoped().First(&kept, deleted.ID).Error)
				require.Equal(t, deleted, kept)
				var savedVotes []entity.Vote
				require.NoError(t, storage.Order("id").Find(&savedVotes).Error)
				expected := append([]entity.Vote(nil), votes...)
				expected[0].TargetID, expected[1].TargetID = pages[0].ID, pages[0].ID
				require.Equal(t, expected, savedVotes)
			}
		})
	}
}
