package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCreateSite(t *testing.T) {
	reloadTestDatabase()

	t.Run("Create New Site", func(t *testing.T) {
		siteName := "TestCreateNewSite"

		result := FindCreateSite(siteName)
		assert.False(t, result.IsEmpty(), "直接获取创建后的站点数据有问题")
		assert.Equal(t, siteName, result.Name)

		findSite := FindSite(siteName)
		assert.False(t, findSite.IsEmpty(), "找不到创建后的站点")
		assert.Equal(t, result.ToCooked(), findSite.ToCooked(), "创建后的站点数据有问题")
	})

	t.Run("Find Existed Site", func(t *testing.T) {
		result := FindCreateSite("Site A")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, "http://localhost:8080/,https://qwqaq.com", result.Urls)
	})
}

func TestFindCreatePage(t *testing.T) {
	reloadTestDatabase()

	t.Run("Create New Page", func(t *testing.T) {
		var (
			pageKey   = "/NewPage.html"
			pageTitle = "New Page Title"
			siteName  = "Site A"
		)

		result := FindCreatePage(pageKey, pageTitle, siteName)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, pageKey, result.Key)
		assert.Equal(t, pageTitle, result.Title)
		assert.Equal(t, siteName, result.SiteName)

		findPage := FindPage(pageKey, siteName)
		assert.False(t, findPage.IsEmpty(), "找不到创建后的页面")
		assert.Equal(t, result.ToCooked(), findPage.ToCooked(), "创建后的页面数据有问题")
	})

	t.Run("Find Existed Page", func(t *testing.T) {
		result := FindCreatePage("/test/1000.html", "", "Site A")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, FindPage("/test/1000.html", "Site A").ToCooked(), result.ToCooked())
	})
}

func TestFindCreateUser(t *testing.T) {
	reloadTestDatabase()

	t.Run("Create New User", func(t *testing.T) {
		var (
			userName  = "NewUser"
			userEmail = "NewUser@gmail.com"
			userLink  = "https://qwqaq.com"
		)

		result := FindCreateUser(userName, userEmail, userLink)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, userName, result.Name)
		assert.Equal(t, userEmail, result.Email)
		assert.Equal(t, userLink, result.Link)

		findUser := FindUser(userName, userEmail)
		assert.False(t, findUser.IsEmpty(), "找不到创建后的用户")
		assert.Equal(t, result.ToCooked(), findUser.ToCooked(), "创建后的用户数据有问题")
	})

	t.Run("Find Existed User", func(t *testing.T) {
		result := FindCreateUser("userA", "user_a@qwqaq.com", "")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, FindUser("userA", "user_a@qwqaq.com").ToCooked(), result.ToCooked())
	})
}
