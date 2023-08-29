package dao

import (
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB

	// please access cache instance by invoking CacheAction func
	cache *DaoCache
}

func NewDao(db *gorm.DB) *Dao {
	dao := &Dao{
		db: db,
	}

	dao.MigrateModels()

	return dao
}

func (dao *Dao) DB() *gorm.DB {
	return dao.db
}

func (dao *Dao) Clone() *Dao {
	clone := *dao

	return &clone
}

func (dao *Dao) SetCache(cache *DaoCache) {
	dao.cache = cache
}

func (dao *Dao) CacheAction(fn func(cache *DaoCache)) {
	if dao.cache != nil {
		fn(dao.cache)
	}
}
