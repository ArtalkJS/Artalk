package dao_test

import (
	"errors"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFailedUpdatesKeepCachedDataConsistentWithStorage(t *testing.T) {
	app, err := test.NewTestApp()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	cache := newTestDaoCache(t)
	t.Cleanup(cache.Close)
	d := app.Dao()
	d.SetCache(cache)

	user := d.FindUserByID(1001)
	site := d.FindSite("Site A")
	page := d.FindCreatePage("/cache-write-failure", "", site.Name)
	comment := entity.Comment{Content: "Original", UserID: user.ID, PageKey: page.Key, SiteName: site.Name}
	require.NoError(t, d.CreateComment(&comment))
	require.NotZero(t, user.ID)
	require.NotZero(t, site.ID)
	require.NotZero(t, page.ID)

	writeErr := errors.New("simulated database write failure")
	require.NoError(t, d.DB().Callback().Update().Before("gorm:update").Register("test:write_failure",
		func(tx *gorm.DB) { tx.AddError(writeErr) }))
	t.Cleanup(func() { require.NoError(t, d.DB().Callback().Update().Remove("test:write_failure")) })

	originalLink, originalURLs, originalTitle := user.Link, site.Urls, page.Title
	user.Link = "https://unsaved.example.org"
	site.Urls = "https://unsaved.example.org"
	page.Title = "Unsaved title"
	comment.Content = "Unsaved content"

	require.ErrorIs(t, d.UpdateUser(&user), writeErr)
	require.Equal(t, originalLink, d.FindUserByID(user.ID).Link)
	require.Equal(t, originalLink, d.FindUser(user.Name, user.Email).Link)

	require.ErrorIs(t, d.UpdateSite(&site), writeErr)
	require.Equal(t, originalURLs, d.FindSite(site.Name).Urls)

	require.ErrorIs(t, d.UpdatePage(&page), writeErr)
	require.Equal(t, originalTitle, d.FindPage(page.Key, site.Name).Title)

	require.ErrorIs(t, d.UpdateComment(&comment), writeErr)
	require.Equal(t, "Original", d.FindComment(comment.ID).Content)
}
