package model

import (
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func SetDB(db *gorm.DB) {
	dbInstance = db
}

func DB() *gorm.DB {
	return dbInstance
}
