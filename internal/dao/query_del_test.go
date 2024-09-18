package dao_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestDelComment(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	comment := app.Dao().FindComment(1000)
	assert.False(t, comment.IsEmpty(), "评论找不到")

	err := app.Dao().DelComment(&comment)
	assert.NoError(t, err, "评论删除错误")

	assert.True(t, app.Dao().FindComment(1000).IsEmpty(), "评论没有删成功")
}

func TestDelCommentChildren(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	parentID := uint(1000)
	err := app.Dao().DelCommentChildren(parentID)
	assert.NoError(t, err, "评论删除错误")

	assert.True(t, app.Dao().FindComment(1004).IsEmpty(), "子评论没有删干净")
}

func TestDelPage(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	page := app.Dao().FindPageByID(1000)
	assert.False(t, page.IsEmpty(), "页面找不到")

	err := app.Dao().DelPage(&page)
	assert.NoError(t, err, "页面删除发生错误")

	assert.True(t, app.Dao().FindPageByID(1000).IsEmpty(), "页面没有删成功")
	assert.True(t, app.Dao().FindComment(1004).IsEmpty(), "页面评论没有删干净")
}

func TestDelSite(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	site := app.Dao().FindSiteByID(1000)
	assert.False(t, site.IsEmpty(), "站点找不到")

	err := app.Dao().DelSite(&site)
	assert.NoError(t, err, "站点删除发生错误")

	assert.True(t, app.Dao().FindSiteByID(1000).IsEmpty(), "站点没有删成功")
	assert.True(t, app.Dao().FindPageByID(1000).IsEmpty(), "站点页面没有删干净")
	assert.True(t, app.Dao().FindComment(1004).IsEmpty(), "站点评论没有删干净")
}

func TestDelUser(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	user := app.Dao().FindUserByID(1000)
	assert.False(t, user.IsEmpty(), "用户找不到")

	err := app.Dao().DelUser(&user)
	assert.NoError(t, err, "用户删除发生错误")

	assert.True(t, app.Dao().FindUserByID(1000).IsEmpty(), "用户没有删成功")

	// Check related records cleaned
	var commentsCount int64
	app.Dao().DB().Where("user_id = ?", 1000).Model(&entity.Comment{}).Count(&commentsCount)
	assert.Equal(t, int64(0), commentsCount, "User comments not cleaned")

	var authIdentityCount int64
	app.Dao().DB().Where("user_id = ?", 1000).Model(&entity.AuthIdentity{}).Count(&authIdentityCount)
	assert.Equal(t, int64(0), authIdentityCount, "User auth identities not cleaned")
}
