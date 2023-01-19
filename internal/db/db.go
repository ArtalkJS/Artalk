package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbInstance *gorm.DB

func SetDB(db *gorm.DB) {
	dbInstance = db
}

func DB() *gorm.DB {
	return dbInstance
}

var gormConfig *gorm.Config

func InitDB() {
	var err error
	db, err := OpenDB(config.Instance.DB.Type, config.Instance.DB.Dsn)
	if err != nil {
		logrus.Error("[DB] ", "Init database error: ", err)
		os.Exit(1)
	}
	SetDB(db)
}

func OpenDB(dbType config.DBType, dsn string) (*gorm.DB, error) {
	dbConf := config.Instance.DB

	gormConfig = &gorm.Config{
		Logger: NewGormLogger(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Instance.DB.TablePrefix,
		},
	}

	if dsn == "" {
		switch dbType {
		case config.TypeSQLite:
			if dbConf.File == "" {
				logrus.Fatal("Please set `db.file` option in config file to specify a sqlite database path")
			}
			dsn = dbConf.File
		case config.TypePostgreSQL:
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				dbConf.Host,
				dbConf.User,
				dbConf.Password,
				dbConf.Name,
				dbConf.Port)
		case config.TypeMySql, config.TypeMSSQL:
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
				dbConf.User,
				dbConf.Password,
				dbConf.Host,
				dbConf.Port,
				dbConf.Name,
				dbConf.Charset,
			)
		}
	}

	switch dbType {
	case config.TypeSQLite:
		return OpenSQLite(dsn)
	case config.TypeMySql:
		return OpenMySql(dsn)
	case config.TypePostgreSQL:
		return OpenPostgreSQL(dsn)
	case config.TypeMSSQL:
		return OpenSqlServer(dsn)
	}

	return nil, errors.New(`unsupported database type "` + string(dbType) + `"`)
}

func OpenSQLite(filename string) (*gorm.DB, error) {
	if err := utils.EnsureDir(filepath.Dir(filename)); err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(filename), gormConfig)
}

func OpenMySql(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), gormConfig)
}

func OpenPostgreSQL(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), gormConfig)
}

func OpenSqlServer(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(dsn), gormConfig)
}
