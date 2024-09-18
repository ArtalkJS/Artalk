package dao_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTestDaoCache(t *testing.T) *dao.DaoCache {
	cache := newTestCache(t)
	return dao.NewCacheAdaptor(cache)
}

func TestNewCacheAdaptor(t *testing.T) {
	cache := newTestDaoCache(t)
	defer cache.Close()

	assert.NotNil(t, cache)
}

func TestCache(t *testing.T) {
	user := entity.User{
		Model: gorm.Model{
			ID: 233,
		},
		Name:     "Artalk",
		Email:    "artalkjs@gmail.com",
		Password: "12345",
	}

	site := entity.Site{
		Model: gorm.Model{
			ID: 123,
		},
		Name: "2333",
	}

	page := entity.Page{
		Model: gorm.Model{
			ID: 343,
		},
		Key:      "/links.html",
		Title:    "Hello World",
		SiteName: site.Name,
	}

	comment := entity.Comment{
		Model: gorm.Model{
			ID: 2321,
		},
		Content:  "abcdefg",
		Rid:      0,
		UserID:   user.ID,
		PageKey:  page.Key,
		SiteName: page.SiteName,
	}

	childComment := entity.Comment{
		Model: gorm.Model{
			ID: 777,
		},
		Content:  "233",
		Rid:      comment.ID,
		UserID:   user.ID,
		PageKey:  page.Key,
		SiteName: page.SiteName,
	}

	tests := []struct {
		name     string
		data     any
		saveFunc func(cache *dao.DaoCache) error
		keys     []string
		delFunc  func(cache *dao.DaoCache) error
	}{
		// User Cache Test
		{
			data: &user,
			saveFunc: func(cache *dao.DaoCache) error {
				return cache.UserCacheSave(&user)
			},
			keys: []string{
				fmt.Sprintf(dao.UserByIDKey, user.ID),
				fmt.Sprintf(dao.UserByNameEmailKey, strings.ToLower(user.Name), strings.ToLower(user.Email)),
			},
			delFunc: func(cache *dao.DaoCache) error {
				cache.UserCacheDel(&user)
				return nil
			},
		},

		// Site Cache Test
		{
			data: &site,
			saveFunc: func(cache *dao.DaoCache) error {
				return cache.SiteCacheSave(&site)
			},
			keys: []string{
				fmt.Sprintf(dao.SiteByIDKey, site.ID),
				fmt.Sprintf(dao.SiteByNameKey, site.Name),
			},
			delFunc: func(cache *dao.DaoCache) error {
				cache.SiteCacheDel(&site)
				return nil
			},
		},

		// Page Cache Test
		{
			data: &page,
			saveFunc: func(cache *dao.DaoCache) error {
				return cache.PageCacheSave(&page)
			},
			keys: []string{
				fmt.Sprintf(dao.PageByIDKey, page.ID),
			},
			delFunc: func(cache *dao.DaoCache) error {
				cache.PageCacheDel(&page)
				return nil
			},
		},

		// Comment Cache Test
		{
			data: &comment,
			saveFunc: func(cache *dao.DaoCache) error {
				return cache.CommentCacheSave(&comment)
			},
			keys: []string{
				fmt.Sprintf(dao.CommentByIDKey, comment.ID),
			},
			delFunc: func(cache *dao.DaoCache) error {
				cache.CommentCacheDel(&comment)
				return nil
			},
		},

		// Child Comment Cache Test
		{
			data: &childComment,
			saveFunc: func(cache *dao.DaoCache) error {
				return cache.CommentCacheSave(&childComment)
			},
			keys: []string{
				fmt.Sprintf(dao.CommentByIDKey, childComment.ID),
			},
			delFunc: func(cache *dao.DaoCache) error {
				cache.CommentCacheDel(&childComment)
				return nil
			},
		},
	}

	for _, tt := range tests {
		for _, key := range tt.keys {
			t.Run(fmt.Sprintf("%T", tt.data)+"/Key="+key, func(t *testing.T) {
				// run test for each key
				cache := newTestDaoCache(t)
				defer cache.Close()

				// func to let Marshaller know the type of the `tt.data`,
				// rather than use `var data any`
				createEmptyTypedData := func() any {
					p := reflect.ValueOf(tt.data).Elem()
					return reflect.New(p.Type()).Interface()
				}

				// Check if key exists before save
				t.Run("FindBeforeSave", func(t *testing.T) {
					var data any
					err := cache.FindCache(key, &data)
					assert.Error(t, err) // Key should not exist before save
					assert.Empty(t, data)
				})

				// save cache
				t.Run("Save", func(t *testing.T) {
					err := tt.saveFunc(cache)
					assert.NoError(t, err)
				})

				// Check if key exists after save
				t.Run("FindAfterSave", func(t *testing.T) {
					data := createEmptyTypedData()

					err := cache.FindCache(key, &data)
					if assert.NoError(t, err) {
						assert.EqualValues(t, tt.data, data) // data should equal after cache save
					}
				})

				// Check if key exists after delete
				t.Run("DeleteAfterSave", func(t *testing.T) {
					err := tt.delFunc(cache)
					assert.NoError(t, err)
				})

				t.Run("FindAfterDelete", func(t *testing.T) {
					var data any
					err := cache.FindCache(key, &data)
					assert.Error(t, err)
					assert.Empty(t, data)
				})
			})
		}
	}
}
