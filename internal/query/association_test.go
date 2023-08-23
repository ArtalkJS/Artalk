package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchUserForComment(t *testing.T) {
	reloadTestDatabase()

	comment := FindComment(1001)
	associatedUser := FetchUserForComment(&comment)
	realUser := FindUserByID(1001)

	assert.False(t, associatedUser.IsEmpty(), "评论关联用户找不到")
	assert.Equal(t, CookUser(&realUser), CookUser(&associatedUser), "评论关联用户数据不一致")
}

func TestFetchPageForComment(t *testing.T) {
	reloadTestDatabase()

	comment := FindComment(1001)
	realPage := FindPage("/test/1000.html", "Site A")
	associatedPage := FetchPageForComment(&comment)

	assert.False(t, associatedPage.IsEmpty(), "评论关联页面找不到")
	assert.Equal(t, CookPage(&realPage), CookPage(&associatedPage), "评论关联页面数据不一致")
}

func TestFetchSiteForComment(t *testing.T) {
	reloadTestDatabase()

	comment := FindComment(1001)
	realSite := FindSite("Site A")
	associatedSite := FetchSiteForComment(&comment)

	assert.False(t, associatedSite.IsEmpty(), "评论关联站点找不到")
	assert.Equal(t, CookSite(&realSite), CookSite(&associatedSite), "评论关联站点数据不一致")
}

func TestFetchSiteForPage(t *testing.T) {
	reloadTestDatabase()

	page := FindPageByID(1000)
	realSite := FindSite("Site A")
	associatedSite := FetchSiteForPage(&page)

	assert.False(t, associatedSite.IsEmpty(), "页面关联站点找不到")
	assert.Equal(t, CookSite(&realSite), CookSite(&associatedSite), "页面关联站点数据不一致")
}

func TestFetchCommentForNotify(t *testing.T) {
	reloadTestDatabase()

	notify := FindNotify(1000, 1001)
	realComment := FindComment(1001)
	associatedComment := FetchCommentForNotify(&notify)

	assert.False(t, associatedComment.IsEmpty(), "通知关联评论找不到")
	assert.Equal(t, CookComment(&realComment), CookComment(&associatedComment), "通知关联评论数据不一致")
}

func TestFetchUserForNotify(t *testing.T) {
	reloadTestDatabase()

	notify := FindNotify(1000, 1001)
	realUser := FindUserByID(1000)
	associatedUser := FetchUserForNotify(&notify)

	assert.False(t, associatedUser.IsEmpty(), "通知关联用户找不到")
	assert.Equal(t, CookUser(&realUser), CookUser(&associatedUser), "通知关联用户数据不一致")
}
