package dao_test

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/test"
	"github.com/stretchr/testify/assert"
)

func TestFindComment(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	type args struct {
		id uint
	}
	type wants struct {
		id        uint
		rid       uint
		user_id   uint
		page_key  string
		site_name string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "评论 ID=1000", args: args{id: 1000}, wants: wants{id: 1000, rid: 0, user_id: 1000, page_key: "/test/1000.html", site_name: "Site A"}},
		{name: "评论 ID=1001", args: args{id: 1001}, wants: wants{id: 1001, rid: 1000, user_id: 1001, page_key: "/test/1000.html", site_name: "Site A"}},
		{name: "评论 ID=1006", args: args{id: 1006}, wants: wants{id: 1006, rid: 0, user_id: 1001, page_key: "/site_b/1001.html", site_name: "Site B"}},
		{name: "不存在的评论", args: args{id: 9999}, wants: wants{id: 0, rid: 0, user_id: 0, page_key: "", site_name: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.Dao().FindComment(tt.args.id)
			assert.Equal(t, tt.wants.id, got.ID)
			assert.Equal(t, tt.wants.rid, got.Rid)
			assert.Equal(t, tt.wants.user_id, got.UserID)
			assert.Equal(t, tt.wants.page_key, got.PageKey)
			assert.Equal(t, tt.wants.site_name, got.SiteName)
			if tt.name != "不存在的评论" {
				assert.NotEmpty(t, got.Content)
			}
		})
	}
}

func TestFindCommentRootID(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	tests := []struct {
		rid    uint
		rootID uint
	}{
		{0, 0},
		{1002, 1000},
	}

	for _, tt := range tests {
		got := app.Dao().FindCommentRootID(tt.rid)
		assert.Equal(t, tt.rootID, got)
	}
}

func TestFindCommentChildrenShallow(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Children Found", func(t *testing.T) {
		result := app.Dao().FindCommentChildrenShallow(1001)
		assert.Equal(t, 2, len(result))
	})

	t.Run("No Children Found", func(t *testing.T) {
		result := app.Dao().FindCommentChildrenShallow(1005)
		assert.Empty(t, result)
	})
}

func TestFindCommentChildren(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Children Found", func(t *testing.T) {
		result := app.Dao().FindCommentChildren(1000)
		assert.Equal(t, 4, len(result))
	})

	t.Run("No Children Found", func(t *testing.T) {
		result := app.Dao().FindCommentChildren(1005)
		assert.Empty(t, result)
	})
}

func TestFindUser(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("User Found", func(t *testing.T) {
		result := app.Dao().FindUser("admin", "admin@qwqaq.com")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, uint(1000), result.ID)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, "admin@qwqaq.com", result.Email)
		assert.Equal(t, "123456", result.Password)
		assert.Equal(t, true, result.IsAdmin)
		assert.Equal(t, "管理员", result.BadgeName)
	})

	t.Run("User not Found", func(t *testing.T) {
		result := app.Dao().FindUser("NoUser", "NoUser@example.org")
		assert.True(t, result.IsEmpty())
	})
}

func TestFindUserByID(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("User Found", func(t *testing.T) {
		result := app.Dao().FindUserByID(1000)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, uint(1000), result.ID)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, "admin@qwqaq.com", result.Email)
		assert.Equal(t, "123456", result.Password)
		assert.Equal(t, true, result.IsAdmin)
		assert.Equal(t, "管理员", result.BadgeName)
	})

	t.Run("User not Found", func(t *testing.T) {
		result := app.Dao().FindUserByID(9999)
		assert.True(t, result.IsEmpty())
	})
}

func TestFindPage(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Page Found", func(t *testing.T) {
		result := app.Dao().FindPage("/site_b/1001.html", "Site B")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, uint(1001), result.ID)
		assert.Equal(t, "测试页面标题 1001", result.Title)
		assert.Equal(t, true, result.AdminOnly)
	})

	t.Run("Page not Found", func(t *testing.T) {
		result := app.Dao().FindPage("/NotExistPage", "NotExistSite")
		assert.True(t, result.IsEmpty())
	})
}

func TestFindPageByID(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Page Found", func(t *testing.T) {
		result := app.Dao().FindPageByID(1001)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, uint(1001), result.ID)
		assert.Equal(t, "测试页面标题 1001", result.Title)
		assert.Equal(t, true, result.AdminOnly)
	})

	t.Run("Page not Found", func(t *testing.T) {
		result := app.Dao().FindPageByID(9999)
		assert.True(t, result.IsEmpty())
	})
}

func TestFindSite(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Site Found", func(t *testing.T) {
		result := app.Dao().FindSite("Site A")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, uint(1000), result.ID)
		assert.Equal(t, "Site A", result.Name)
		assert.Equal(t, "http://localhost:8080/,https://qwqaq.com", result.Urls)
	})

	t.Run("Site not Found", func(t *testing.T) {
		result := app.Dao().FindSite("NotExistSite")
		assert.True(t, result.IsEmpty())
	})
}

func TestFindSiteByID(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Site Found", func(t *testing.T) {
		result := app.Dao().FindSiteByID(1000)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, uint(1000), result.ID)
		assert.Equal(t, "Site A", result.Name)
		assert.Equal(t, "http://localhost:8080/,https://qwqaq.com", result.Urls)
	})

	t.Run("Site not Found", func(t *testing.T) {
		result := app.Dao().FindSiteByID(9999)
		assert.True(t, result.IsEmpty())
	})
}

func TestFindAllSites(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	allSites := app.Dao().FindAllSites()
	assert.GreaterOrEqual(t, len(allSites), 1)
}

func TestGetAllAdmins(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	allAdmins := app.Dao().GetAllAdmins()
	assert.GreaterOrEqual(t, len(allAdmins), 1, "GetAllAdmins() not works")

	t.Run("Test modify and get admins", func(t *testing.T) {
		getContainsUser := func(userID uint) bool {
			for _, a := range app.Dao().GetAllAdmins() {
				if a.ID == userID {
					return true
				}
			}
			return false
		}

		// create
		var adminID uint = 0
		admin := entity.User{
			Name:    "TestAdmin",
			Email:   "admin@test.com",
			IsAdmin: true,
		}
		app.Dao().CreateUser(&admin)
		adminID = admin.ID
		assert.NotZero(t, adminID, "user create failed")

		assert.Equal(t, true, getContainsUser(adminID), "admin not found after create")

		// update
		admin.IsAdmin = false
		app.Dao().UpdateUser(&admin)
		assert.Equal(t, false, getContainsUser(adminID), "admin still exists after update")

		// re-update
		admin.IsAdmin = true
		app.Dao().UpdateUser(&admin)
		assert.Equal(t, true, getContainsUser(adminID), "admin not found after re-update to admin")

		// delete
		app.Dao().DelUser(&admin)

		// not contains admin
		assert.Equal(t, false, getContainsUser(adminID), "admin still exists after delete")
	})
}

func TestIsAdminUser(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	assert.Equal(t, true, app.Dao().IsAdminUser(1000))
}

func TestIsAdminUserByNameEmail(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	assert.Equal(t, true, app.Dao().IsAdminUserByNameEmail("admin", "admin@qwqaq.com"))
}
