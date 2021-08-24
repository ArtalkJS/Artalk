package lib

import (
	"github.com/ArtalkJS/Artalk-API-Go/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDb(dbFile string) (err error) {
	DB, err = gorm.Open(sqlite.Open(config.Instance.DB.Dsn), &gorm.Config{})
	return
}
