package dao_test

import (
	"testing"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTestDao(t *testing.T) (*dao.Dao, *gorm.DB) {
	db, err := db.OpenSQLite("file::memory:?cache=shared", &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	dao := dao.NewDao(db)
	return dao, db
}

func TestBase(t *testing.T) {
	d, _ := newTestDao(t)
	assert.NotNil(t, d)
}

func TestDB(t *testing.T) {
	d, db := newTestDao(t)
	assert.Equal(t, db, d.DB())
}

func newTestCache(t *testing.T) *cache.Cache {
	c, err := cache.New(config.CacheConf{
		Type: config.CacheTypeBuiltin,
	})
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestDaoCache(t *testing.T) {
	d, _ := newTestDao(t)

	cacheAdaptor := dao.NewCacheAdaptor(newTestCache(t))
	defer cacheAdaptor.Close()

	t.Run("SetCache", func(t *testing.T) {
		d.SetCache(cacheAdaptor)
	})

	t.Run("CacheAction_Valid", func(t *testing.T) {
		cacheActionFnInvoked := false
		d.CacheAction(func(cache *dao.DaoCache) {
			cacheActionFnInvoked = true
			assert.Equal(t, cacheAdaptor, cache)
		})
		assert.True(t, cacheActionFnInvoked)
	})

	t.Run("CacheAction_Invalid", func(t *testing.T) {
		d, _ := newTestDao(t)
		cacheActionFnInvoked := false
		d.CacheAction(func(cache *dao.DaoCache) {
			cacheActionFnInvoked = true
		})
		assert.False(t, cacheActionFnInvoked)
	})
}
