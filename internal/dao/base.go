package dao

import (
	"github.com/ArtalkJS/Artalk/internal/dao/dao_cache"
	"gorm.io/gorm"
)

type Dao struct {
	db    *gorm.DB
	cache *dao_cache.Cache
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

func (dao *Dao) SetCache(cache *dao_cache.Cache) {
	dao.cache = cache
}

func (dao *Dao) Cache() *dao_cache.Cache {
	return dao.cache
}
