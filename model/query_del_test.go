package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelComment(t *testing.T) {
	reloadTestDatabase()

	comment := FindComment(1000)
	assert.False(t, comment.IsEmpty(), "评论找不到")

	err := DelComment(&comment)
	assert.NoError(t, err, "评论删除错误")

	assert.True(t, FindComment(1000).IsEmpty(), "评论没有删成功")
}

func TestDelCommentChildren(t *testing.T) {
	reloadTestDatabase()

	parentID := uint(1000)
	err := DelCommentChildren(parentID)
	assert.NoError(t, err, "评论删除错误")

	assert.True(t, FindComment(1004).IsEmpty(), "子评论没有删干净")
}

func TestDelPage(t *testing.T) {
	reloadTestDatabase()

	page := FindPageByID(1000)
	assert.False(t, page.IsEmpty(), "页面找不到")

	err := DelPage(&page)
	assert.NoError(t, err, "页面删除发生错误")

	assert.True(t, FindPageByID(1000).IsEmpty(), "页面没有删成功")
	assert.True(t, FindComment(1004).IsEmpty(), "页面评论没有删干净")
}

func TestDelSite(t *testing.T) {
	reloadTestDatabase()

	site := FindSiteByID(1000)
	assert.False(t, site.IsEmpty(), "站点找不到")

	err := DelSite(&site)
	assert.NoError(t, err, "站点删除发生错误")

	assert.True(t, FindSiteByID(1000).IsEmpty(), "站点没有删成功")
	assert.True(t, FindPageByID(1000).IsEmpty(), "站点页面没有删干净")
	assert.True(t, FindComment(1004).IsEmpty(), "站点评论没有删干净")
}

func TestDelUser(t *testing.T) {
	reloadTestDatabase()

	user := FindUserByID(1000)
	assert.False(t, user.IsEmpty(), "用户找不到")

	err := DelUser(&user)
	assert.NoError(t, err, "用户删除发生错误")

	assert.True(t, FindUserByID(1000).IsEmpty(), "用户没有删成功")
}
