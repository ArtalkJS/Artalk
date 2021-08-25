package lib

import (
	"path/filepath"

	"github.com/ArtalkJS/ArtalkGo/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

var gormConfig = &gorm.Config{
	Logger: NewGormLogger(),
}

func OpenDB() (err error) {
	switch config.Instance.DB.Type {
	case config.TypeMySql:
		err = OpenMySqlDb()
	case config.TypeSqlite:
		err = OpenSqliteDb()
	}
	return
}

func OpenMySqlDb() (err error) {
	DB, err = gorm.Open(mysql.Open(config.Instance.DB.Dsn), gormConfig)
	return
}

func OpenSqliteDb() (err error) {
	filename := config.Instance.DB.Dsn
	if err := EnsureDir(filepath.Dir(filename)); err != nil {
		return err
	}

	DB, err = gorm.Open(sqlite.Open(filename), gormConfig)
	return
}
