package cache

import (
	"github.com/ArtalkJS/ArtalkGo/internal/db"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.DB()
}
