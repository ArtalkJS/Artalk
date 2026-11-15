package dao_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestCacheWarmUpAndFlush(t *testing.T) {
	d, storage := newTestDao(t)
	t.Cleanup(func() { require.NoError(t, db.CloseDB(storage)) })

	user := entity.User{Name: "Reader", Email: "warmup@example.org"}
	site := entity.Site{Name: "Cached site"}
	page := entity.Page{Key: "/cached", SiteName: site.Name, Title: "Cached page"}
	for _, record := range []any{&user, &site, &page} {
		require.NoError(t, storage.Create(record).Error)
	}
	parent := entity.Comment{Content: "Parent", SiteName: site.Name, PageKey: page.Key, UserID: user.ID}
	require.NoError(t, storage.Create(&parent).Error)
	reply := entity.Comment{Content: "Reply", SiteName: site.Name, PageKey: page.Key, UserID: user.ID, Rid: parent.ID}
	require.NoError(t, storage.Create(&reply).Error)

	c := newTestDaoCache(t)
	t.Cleanup(c.Close)
	d.SetCache(c)
	d.CacheWarmUp()

	// Change the source directly to distinguish warmed results from database reads.
	require.NoError(t, storage.Model(&user).Update("link", "https://updated.example.org").Error)
	require.NoError(t, storage.Model(&site).Update("urls", "https://updated.example.org").Error)
	require.NoError(t, storage.Model(&page).Update("title", "Updated page").Error)
	require.NoError(t, storage.Model(&reply).Update("content", "Updated reply").Error)
	another := entity.Comment{Content: "New reply", SiteName: site.Name, PageKey: page.Key, UserID: user.ID, Rid: parent.ID}
	require.NoError(t, storage.Create(&another).Error)

	require.Empty(t, d.FindUser(user.Name, user.Email).Link)
	require.Empty(t, d.FindSite(site.Name).Urls)
	require.Equal(t, "Cached page", d.FindPage(page.Key, site.Name).Title)
	children := d.FindCommentChildrenShallow(parent.ID)
	require.Len(t, children, 1)
	require.Equal(t, "Reply", children[0].Content)

	d.CacheFlushAll()
	require.Equal(t, "https://updated.example.org", d.FindUser(user.Name, user.Email).Link)
	require.Equal(t, "https://updated.example.org", d.FindSite(site.Name).Urls)
	require.Equal(t, "Updated page", d.FindPage(page.Key, site.Name).Title)
	children = d.FindCommentChildrenShallow(parent.ID)
	require.Len(t, children, 2)
	require.Equal(t, "Updated reply", d.FindComment(reply.ID).Content)
}
