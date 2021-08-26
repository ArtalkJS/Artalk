package lib

import (
	"path/filepath"

	"github.com/ArtalkJS/ArtalkGo/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

var gormConfig = &gorm.Config{
	Logger: NewGormLogger(),
}

func OpenDB() (err error) {
	switch config.Instance.DB.Type {
	case config.TypeMySql:
		err = OpenMySql()
	case config.TypeSQLite:
		err = OpenSQLite()
	case config.TypePostgreSQL:
		err = OpenPostgreSQL()
	case config.TypeSqlServer:
		err = OpenSqlServer()
	}
	return
}

func OpenMySql() (err error) {
	DB, err = gorm.Open(mysql.Open(config.Instance.DB.Dsn), gormConfig)
	return
}

func OpenSQLite() (err error) {
	filename := config.Instance.DB.Dsn
	if err := EnsureDir(filepath.Dir(filename)); err != nil {
		return err
	}

	DB, err = gorm.Open(sqlite.Open(filename), gormConfig)
	return
}

func OpenPostgreSQL() (err error) {
	DB, err = gorm.Open(postgres.Open(config.Instance.DB.Dsn), &gorm.Config{})
	return
}

func OpenSqlServer() (err error) {
	DB, err = gorm.Open(sqlserver.Open(config.Instance.DB.Dsn), &gorm.Config{})
	return
}
