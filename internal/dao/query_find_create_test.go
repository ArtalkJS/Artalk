package dao_test

import (
	"testing"

	"github.com/ArtalkJS/Artalk/test"
	"github.com/stretchr/testify/assert"
)

func TestFindCreateSite(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Create New Site", func(t *testing.T) {
		siteName := "TestCreateNewSite"

		result := app.Dao().FindCreateSite(siteName)
		assert.False(t, result.IsEmpty(), "直接获取创建后的站点数据有问题")
		assert.Equal(t, siteName, result.Name)

		findSite := app.Dao().FindSite(siteName)
		assert.False(t, findSite.IsEmpty(), "找不到创建后的站点")
		assert.Equal(t, app.Dao().CookSite(&result), app.Dao().CookSite(&findSite), "创建后的站点数据有问题")
	})

	t.Run("Find Existed Site", func(t *testing.T) {
		result := app.Dao().FindCreateSite("Site A")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, "http://localhost:8080/,https://qwqaq.com", result.Urls)
	})
}

func TestFindCreatePage(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Create New Page", func(t *testing.T) {
		var (
			pageKey   = "/NewPage.html"
			pageTitle = "New Page Title"
			siteName  = "Site A"
		)

		result := app.Dao().FindCreatePage(pageKey, pageTitle, siteName)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, pageKey, result.Key)
		assert.Equal(t, pageTitle, result.Title)
		assert.Equal(t, siteName, result.SiteName)

		findPage := app.Dao().FindPage(pageKey, siteName)
		assert.False(t, findPage.IsEmpty(), "找不到创建后的页面")
		assert.Equal(t, app.Dao().CookPage(&result), app.Dao().CookPage(&findPage), "创建后的页面数据有问题")
	})

	t.Run("Find Existed Page", func(t *testing.T) {
		result := app.Dao().FindCreatePage("/test/1000.html", "", "Site A")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, app.Dao().FindPage("/test/1000.html", "Site A"), result)
	})
}

func TestFindCreateUser(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Create New User", func(t *testing.T) {
		var (
			userName  = "NewUser"
			userEmail = "NewUser@gmail.com"
			userLink  = "https://qwqaq.com"
		)

		result, err := app.Dao().FindCreateUser(userName, userEmail, userLink)
		assert.NoError(t, err)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, userName, result.Name)
		assert.Equal(t, userEmail, result.Email)
		assert.Equal(t, userLink, result.Link)

		findUser := app.Dao().FindUser(userName, userEmail)
		assert.False(t, findUser.IsEmpty(), "找不到创建后的用户")
		assert.Equal(t, app.Dao().CookUser(&result), app.Dao().CookUser(&findUser), "创建后的用户数据有问题")
	})

	t.Run("Valid User Values", func(t *testing.T) {
		args := []struct {
			name   string
			email  string
			link   string
			result bool
		}{
			{"", "", "", false},
			{"userA", "", "", false},
			{"", "user_a@example.com", "", false},
			{"userB", "user_b", "", false},
			{"userC", "user_c@example.com", "https://xxxx.com", true},
		}
		for _, arg := range args {
			_, err := app.Dao().FindCreateUser(arg.name, arg.email, arg.link)
			assert.Equal(t, arg.result, err == nil, "FindCreateUser(%s, %s, %s) should return %v", arg.name, arg.email, arg.link, arg.result)
		}

		// Invalid user link
		u, err := app.Dao().FindCreateUser("userD", "user_d@example.com", "xxxx.com")
		assert.NoError(t, err)
		assert.Equal(t, "", u.Link, "The user should be create but link is empty because it's invalid")
	})

	t.Run("Find Existed User", func(t *testing.T) {
		result, err := app.Dao().FindCreateUser("userA", "user_a@qwqaq.com", "")
		assert.NoError(t, err)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, app.Dao().FindUser("userA", "user_a@qwqaq.com"), result)
	})
}
