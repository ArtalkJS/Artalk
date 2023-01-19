package cache

import (
	"github.com/ArtalkJS/Artalk/internal/db"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.DB()
}
