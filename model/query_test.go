package model

import (
	"log"
	"os"
	"testing"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func CreateTestData() {
	// 创建站点
	site := Site{
		Model: gorm.Model{
			ID: 100,
		},
		Name: "测试站点",
		Urls: "https://localhost:23360,https://qwqaq.com",
	}
	lib.DB.Create(&site)

	// 测试页面
	page := Page{
		Model: gorm.Model{
			ID: 200,
		},
		Key:      "/test-pages/1.html",
		Title:    "测试页面 1",
		SiteName: site.Name,
	}
	lib.DB.Create(&page)

	// 测试用户
	normalUser := User{
		Model: gorm.Model{
			ID: 300,
		},
		Name:    "测试普通用户",
		Email:   "test@qwqaq.com",
		Link:    "https://qwqaq.com",
		IsAdmin: false,
	}
	adminUser := User{
		Model: gorm.Model{
			ID: 301,
		},
		Name:    "测试管理员用户",
		Email:   "admin@qwqaq.com",
		Link:    "https://qwqaq.com",
		IsAdmin: true,
	}
	lib.DB.Create(&normalUser)
	lib.DB.Create(&adminUser)

	// 普通用户评论
	comment := Comment{
		Model: gorm.Model{
			ID: 400,
		},
		Content:  "## 测试内容...\n这是一条来自普通用户的评论",
		PageKey:  page.Key,
		SiteName: site.Name,
		UserID:   normalUser.ID,
		Rid:      0,
	}
	lib.DB.Create(&comment)

	// 管理员用户评论
	adminComment := Comment{
		Model: gorm.Model{
			ID: 401,
		},
		Content:  "## 来自管理员的评论\ncheck check check",
		PageKey:  page.Key,
		SiteName: site.Name,
		UserID:   adminUser.ID,
		Rid:      comment.Rid,
	}
	lib.DB.Create(&adminComment)
}

// TestMain 初始化 Mock
func TestMain(m *testing.M) {
	if err := InitTestDB(); err != nil {
		log.Panic("测试数据库初始化发生错误 ", err)
	}

	CreateTestData()

	code := m.Run()
	os.Exit(code)
}

func TestFindComment(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		c := FindComment(400)
		if c.IsEmpty() {
			t.Fatal("comment cannot found")
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		c := FindComment(999)
		if !c.IsEmpty() {
			t.Fatal("comment empty expected")
		}
	})
}

func TestFindCommentScopes(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		c := FindCommentScopes(400, func(db *gorm.DB) *gorm.DB {
			return db.Where(&Comment{
				SiteName: "测试站点",
			})
		})
		if c.IsEmpty() {
			t.Fatal("comment with scope cannot found")
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		c := FindCommentScopes(400, func(db *gorm.DB) *gorm.DB {
			return db.Where(&Comment{
				SiteName: "不存在的测试站点",
			})
		})
		if !c.IsEmpty() {
			t.Fatal("comment with scope empty expected")
		}
	})
}

func TestFindUser(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		u1 := FindUser("测试普通用户", "test@qwqaq.com")
		if u1.IsEmpty() {
			t.Fatal("user 测试普通用户 cannot found")
		}

		u2 := FindUser("测试管理员用户", "admin@qwqaq.com")
		if u2.IsEmpty() {
			t.Fatal("user 测试管理员用户 cannot found")
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		u1 := FindUser("不存在的用户名", "admin@qwqaq.com")
		if !u1.IsEmpty() {
			t.Fatal("user 不存在的用户名 empty expected")
		}

		u2 := FindUser("测试管理员用户", "不存在的邮箱@qwqaq.com")
		if !u2.IsEmpty() {
			t.Fatal("user 不存在的邮箱 empty expected")
		}
	})
}

func TestFindUserByID(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		u1 := FindUserByID(300)
		if u1.IsEmpty() {
			t.Fatal("user 测试普通用户 cannot found")
		}

		u2 := FindUserByID(301)
		if u2.IsEmpty() {
			t.Fatal("user 测试管理员用户 cannot found")
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		u := FindUserByID(999)
		if !u.IsEmpty() {
			t.Fatal("user empty expected")
		}
	})
}

func TestUpdateComment(t *testing.T) {
	comment := Comment{
		Model: gorm.Model{
			ID: 408,
		},
		Content:  "测试内容",
		PageKey:  "test",
		SiteName: "test",
		UserID:   999,
		Rid:      0,
	}
	lib.DB.Create(&comment)

	findComment := FindComment(408)

	modifiedConent := "测试修改后的内容"
	findComment.Content = modifiedConent

	UpdateComment(&findComment)

	findComment2 := FindComment(408)

	if !assert.Equal(t, findComment2.Content, modifiedConent) {
		t.Fatal("修改无效")
	}
}

func TestFindSite(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		s := FindSite("测试站点")
		if s.IsEmpty() {
			t.Fatal("site 测试站点 cannot found")
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		s := FindSite("不存在的站点")
		if !s.IsEmpty() {
			t.Fatal("site 不存在的站点 empty expected")
		}
	})
}

func TestFindSiteByID(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		s := FindSiteByID(100)
		if s.IsEmpty() {
			t.Fatal("site cannot found")
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		s := FindSiteByID(999)
		if !s.IsEmpty() {
			t.Fatal("site empty expected")
		}
	})
}

func TestFindCreateSite(t *testing.T) {
	s := FindCreateSite("测试站点")
	if s.IsEmpty() {
		t.Fatal("已存在的站点找不到")
	}

	newSiteName := "测试站点 FindCreateSite"
	s2 := FindCreateSite(newSiteName)
	if s2.IsEmpty() {
		t.Fatal("创建不存在的站点失败")
	}

	if FindSite(newSiteName).IsEmpty() {
		t.Fatal("创建不存在的站点找不到")
	}
}

func TestNewSite(t *testing.T) {
	siteName := "新增站点"
	siteUrls := "https://qwqaq.com"
	site := NewSite(siteName, siteUrls)

	findSite := FindSite(siteName)
	if findSite.IsEmpty() {
		t.Fatal("创建不存在的站点找不到")
	}
	assert.Equal(t, findSite.Name, site.Name)
	assert.Equal(t, findSite.Urls, site.Urls)
}

func TestUpdateSite(t *testing.T) {
	siteName := "新增站点2"
	siteNameAfter := "新增站点2 修改后"
	siteUrls := "https://qwqaq.com"
	siteUrlsAfter := "https://qwqaq.com/233"
	NewSite(siteName, siteUrls)

	findSite := FindSite(siteName)
	findSite.Name = siteNameAfter
	findSite.Urls = siteUrlsAfter
	UpdateSite(&findSite)

	findSite2 := FindSite(siteNameAfter)
	if findSite2.IsEmpty() {
		t.Fatal("找不到修改后的站点")
	}

	assert.Equal(t, findSite2.Name, siteNameAfter)
	assert.Equal(t, findSite2.Urls, siteUrlsAfter)
}

func TestFindCreatePage(t *testing.T) {
	p := FindCreatePage("/test-pages/1.html", "测试页面 1", "测试站点")
	if p.IsEmpty() {
		t.Fatal("已存在的页面找不到")
	}

	newPageKey := "测试页面 FindCreatePage"
	newPageTitle := "测试页面 FindCreatePage"
	p2 := FindCreatePage(newPageKey, newPageTitle, "测试站点")
	if p2.IsEmpty() {
		t.Fatal("创建不存在的页面失败")
	}

	findPage := FindPage(newPageKey, "测试站点")
	if findPage.IsEmpty() {
		t.Fatal("创建不存在的页面找不到")
	}

	assert.Equal(t, findPage.Key, newPageKey)
	assert.Equal(t, findPage.Title, newPageTitle)
}
