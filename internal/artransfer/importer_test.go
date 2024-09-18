package artransfer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_importArtrans(t *testing.T) {
	t.Run("Import with TargetSiteName, TargetSiteURL and URLResolver", func(t *testing.T) {
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		TEST_SITE_NAME := "test_site"
		TEST_CONTENT := "TestContent"
		TEST_NEGATIVE_SITE := "NegativeSite"
		TEST_SITE_URL := "https://example.com"
		TEST_PAGE_KEY := "/test_page_key.html"
		TEST_USER_NAME := "TestUser"
		TEST_USER_EMAIL := "test_user@exmaple.com"
		TEST_USER_LINK := "https://user_link_example.com"

		params := ImportParams{
			Assumeyes:      true,
			TargetSiteName: TEST_SITE_NAME,
			TargetSiteURL:  TEST_SITE_URL,
			URLResolver:    true,
			URLKeepDomain:  false,
		}

		err := importArtrans(dao.DB(), &params, []*entity.Artran{
			{
				Content:  TEST_CONTENT,
				SiteName: TEST_NEGATIVE_SITE,
				PageKey:  TEST_PAGE_KEY,
				Nick:     TEST_USER_NAME,
				Email:    TEST_USER_EMAIL,
				Link:     TEST_USER_LINK,
			},
		})
		assert.Nil(t, err, "Import should be successful")

		t.Run("Check if the site is created properly", func(t *testing.T) {
			site := findSite(dao.DB(), TEST_SITE_NAME)
			assert.False(t, site.IsEmpty(), "Site should be created")
			assert.Equal(t, TEST_SITE_NAME, site.Name, "Site.Name should be the target site name")
			assert.Equal(t, TEST_SITE_URL, site.Urls, "Site.Urls should be the same")
		})

		t.Run("Check if the user is created properly", func(t *testing.T) {
			var user entity.User
			dao.DB().First(&user)
			assert.False(t, user.IsEmpty(), "User should be created")
			assert.Equal(t, TEST_USER_NAME, user.Name, "User.Name should be the same")
			assert.Equal(t, TEST_USER_EMAIL, user.Email, "User.Email should be the same")
			assert.Equal(t, TEST_USER_LINK, user.Link, "User.Link should be the same")
		})

		t.Run("Check if the comment is created properly", func(t *testing.T) {
			var user entity.User
			dao.DB().First(&user)
			assert.False(t, user.IsEmpty(), "User should be created")

			var comment entity.Comment
			dao.DB().First(&comment)
			assert.Equal(t, TEST_CONTENT, comment.Content)
			assert.Equal(t, TEST_SITE_NAME, comment.SiteName, "SiteName should be the target site name")
			assert.True(t, strings.HasPrefix(comment.PageKey, TEST_SITE_URL), "PageKey should be started with the target site url, since URLResolver is enabled")
			assert.Equal(t, comment.PageKey, TEST_SITE_URL+TEST_PAGE_KEY, "PageKey should be the same")
			assert.Equal(t, user.ID, comment.UserID, "UserID should be the same")
		})

		// Check if the negative site is not created
		negativeSite := findSite(dao.DB(), TEST_NEGATIVE_SITE)
		assert.True(t, negativeSite.IsEmpty(), "Negative site should not be created, The `target_site_name` should be used overwrite the `site_name` of each artran")
	})

	t.Run("Import with StripDomain", func(t *testing.T) {
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		params := ImportParams{
			Assumeyes:     true,
			URLKeepDomain: false,
		}

		err := importArtrans(dao.DB(), &params, []*entity.Artran{
			{
				Content:  "TestContent",
				SiteName: "test_site",
				PageKey:  "https://example.com/test_page_key.html",
			},
		})
		assert.Nil(t, err, "Import should be successful")

		var comment entity.Comment
		dao.DB().First(&comment)
		assert.Equal(t, "/test_page_key.html", comment.PageKey, "PageKey should be the same")
	})

	t.Run("Import with Append SiteURL from Artran", func(t *testing.T) {
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		params := ImportParams{
			Assumeyes: true,
		}

		// Create a site
		dao.DB().Create(&entity.Site{
			Name: "test_site",
			Urls: "https://example.com",
		})

		// Perform import
		err := importArtrans(dao.DB(), &params, []*entity.Artran{
			{
				Content:  "TestContent",
				PageKey:  "/test_page_key.html",
				SiteName: "test_site",
				SiteURLs: "https://example.com,https://example2.com,https://example3.com",
			},
		})
		assert.Nil(t, err, "Import should be successful")

		// Assert
		var site entity.Site
		dao.DB().First(&site)
		assert.Equal(t, "https://example2.com,https://example3.com,https://example.com", site.Urls, "SiteURLs should be appended")
	})

	t.Run("Import with Votes", func(t *testing.T) {
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		params := ImportParams{
			Assumeyes: true,
		}

		// Perform import
		err := importArtrans(dao.DB(), &params, []*entity.Artran{
			{
				Content:  "TestContent",
				PageKey:  "/test_page_key.html",
				SiteName: "test_site",
				VoteUp:   "5",
				VoteDown: "3",
			},
		})
		assert.Nil(t, err, "Import should be successful")

		// Assert
		var comment entity.Comment
		dao.DB().First(&comment)
		assert.Equal(t, 5, comment.VoteUp, "VoteUp should be the same")
		assert.Equal(t, 3, comment.VoteDown, "VoteDown should be the same")

		var upCount int64
		dao.DB().Model(&entity.Vote{}).Where(&entity.Vote{TargetID: comment.ID, Type: entity.VoteTypeCommentUp}).Count(&upCount)
		assert.Equal(t, int64(5), upCount, "VoteUp should be the same")

		var downCount int64
		dao.DB().Model(&entity.Vote{}).Where(&entity.Vote{TargetID: comment.ID, Type: entity.VoteTypeCommentDown}).Count(&downCount)
		assert.Equal(t, int64(3), downCount, "VoteDown should be the same")
	})

	t.Run("Import with DB transaction Rollback", func(t *testing.T) {
		ddb, _ := db.NewTestDB()
		defer db.CloseDB(ddb)
		dao := dao.NewDao(ddb)

		params := ImportParams{
			Assumeyes: true,
		}

		artransJSON := `[
			{"content": "TestContent", "page_key": "/test_page_key.html", "site_name": "test_site", "nick": "TestUser", "email": "abc@example.com"},
			{"content": "TestContent", "page_key": "/test_page_key.html", "site_name": "test_site"},
		]` // without site_name

		// Perform import
		params.JsonData = artransJSON

		// Mock db error
		// make comments table unique on content
		dao.DB().Exec("CREATE UNIQUE INDEX idx_comments_content ON " + dao.GetTableName(&entity.Comment{}) + " (content)")

		err := RunImportArtrans(dao, &params)
		assert.Error(t, err, "Import should be failed")

		// Assert
		entities := []any{entity.Comment{}, entity.Site{}, entity.Page{}, entity.User{}}
		for _, e := range entities {
			var count int64
			err := dao.DB().Model(&e).Count(&count).Error
			assert.Nil(t, err, "Should not be error")
			assert.Equal(t, int64(0), count, fmt.Sprintf("No record should be created for %T", e))
		}
	})
}

func Test_buildRid2RootGenIdMap(t *testing.T) {
	comments := []*entity.Artran{
		{ID: "a", Rid: "0"},

		{ID: "b", Rid: "a"},
		{ID: "c", Rid: "b"},
		{ID: "d", Rid: "c"},
		{ID: "e", Rid: "d"},

		{ID: "f", Rid: "c"},
		{ID: "g", Rid: ""},
		{ID: "h", Rid: "g"},
	}

	genIdMap := buildGenIdMap(comments)
	rootIdMap := buildRid2RootGenIdMap(comments, genIdMap)

	expected := map[string]uint{"": 0, "0": 0, "a": 1, "b": 1, "c": 1, "d": 1, "g": 7}
	assert.Equal(t, expected, rootIdMap)
}
