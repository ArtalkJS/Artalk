package dao_test

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTestDao(t *testing.T) (*dao.Dao, *gorm.DB) {
	db, err := db.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

	dao := dao.NewDao(db)
	return dao, db
}

func TestBase(t *testing.T) {
	dao, ddb := newTestDao(t)
	defer db.CloseDB(ddb)

	assert.NotNil(t, dao)
}

func TestDB(t *testing.T) {
	dao, ddb := newTestDao(t)
	defer db.CloseDB(ddb)

	assert.Equal(t, ddb, dao.DB())
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
