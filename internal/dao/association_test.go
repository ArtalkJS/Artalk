package dao_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestFetchUserForComment(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	comment := app.Dao().FindComment(1001)
	associatedUser := app.Dao().FetchUserForComment(&comment)
	realUser := app.Dao().FindUserByID(1001)

	assert.False(t, associatedUser.IsEmpty(), "评论关联用户找不到")
	assert.Equal(t, app.Dao().CookUser(&realUser), app.Dao().CookUser(&associatedUser), "评论关联用户数据不一致")
}

func TestFetchPageForComment(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	comment := app.Dao().FindComment(1001)
	realPage := app.Dao().FindPage("/test/1000.html", "Site A")
	associatedPage := app.Dao().FetchPageForComment(&comment)

	assert.False(t, associatedPage.IsEmpty(), "评论关联页面找不到")
	assert.Equal(t, app.Dao().CookPage(&realPage), app.Dao().CookPage(&associatedPage), "评论关联页面数据不一致")
}

func TestFetchSiteForComment(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	comment := app.Dao().FindComment(1001)
	realSite := app.Dao().FindSite("Site A")
	associatedSite := app.Dao().FetchSiteForComment(&comment)

	assert.False(t, associatedSite.IsEmpty(), "评论关联站点找不到")
	assert.Equal(t, app.Dao().CookSite(&realSite), app.Dao().CookSite(&associatedSite), "评论关联站点数据不一致")
}

func TestFetchSiteForPage(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	page := app.Dao().FindPageByID(1000)
	realSite := app.Dao().FindSite("Site A")
	associatedSite := app.Dao().FetchSiteForPage(&page)

	assert.False(t, associatedSite.IsEmpty(), "页面关联站点找不到")
	assert.Equal(t, app.Dao().CookSite(&realSite), app.Dao().CookSite(&associatedSite), "页面关联站点数据不一致")
}

func TestFetchCommentForNotify(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	notify := app.Dao().FindNotify(1000, 1001)
	realComment := app.Dao().FindComment(1001)
	associatedComment := app.Dao().FetchCommentForNotify(&notify)

	assert.False(t, associatedComment.IsEmpty(), "通知关联评论找不到")
	assert.Equal(t, app.Dao().CookComment(&realComment), app.Dao().CookComment(&associatedComment), "通知关联评论数据不一致")
}

func TestFetchUserForNotify(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	notify := app.Dao().FindNotify(1000, 1001)
	realUser := app.Dao().FindUserByID(1000)
	associatedUser := app.Dao().FetchUserForNotify(&notify)

	assert.False(t, associatedUser.IsEmpty(), "通知关联用户找不到")
	assert.Equal(t, app.Dao().CookUser(&realUser), app.Dao().CookUser(&associatedUser), "通知关联用户数据不一致")
}
