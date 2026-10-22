package dao_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
)

func TestPagePVPreservesNewerValues(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()
	cache := newTestDaoCache(t)
	defer cache.Close()
	app.Dao().SetCache(cache)
	page := app.Dao().FindCreatePage("/pv-stale", "", "Site A")
	other := page
	require.NoError(t, app.Dao().DB().Model(&entity.Page{}).Where("id = ?", page.ID).
		Updates(map[string]any{"pv": 123, "title": "Admin title", "admin_only": true}).Error)
	// Independent requests may hold the same stale snapshot.
	require.NoError(t, app.Dao().IncrementPagePV(&page))
	require.NoError(t, app.Dao().IncrementPagePV(&other))
	stored := app.Dao().FindPage(page.Key, page.SiteName)
	require.Equal(t, 125, stored.PV)
	require.Equal(t, 125, other.PV)
	require.Equal(t, "Admin title", stored.Title)
	require.True(t, stored.AdminOnly)

	// The historical PV column is nullable; old imports may contain SQL NULL.
	require.NoError(t, app.Dao().DB().Model(&entity.Page{}).Where("id = ?", page.ID).Update("pv", nil).Error)
	require.NoError(t, app.Dao().IncrementPagePV(&page))
	require.Equal(t, 1, page.PV)

	require.NoError(t, app.Dao().DB().Unscoped().Delete(&entity.Page{}, page.ID).Error)
	require.Error(t, app.Dao().IncrementPagePV(&page))
	var count int64
	require.NoError(t, app.Dao().DB().Unscoped().Model(&entity.Page{}).Where("id = ?", page.ID).Count(&count).Error)
	require.Zero(t, count)
}
