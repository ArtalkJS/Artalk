package dao

import (
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func NewDao(db *gorm.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (dao *Dao) DB() *gorm.DB {
	return dao.db
}

func (dao *Dao) Clone() *Dao {
	clone := *dao

	return &clone
}
