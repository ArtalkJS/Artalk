package lib

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

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

func OpenDB(dbType config.DBType, dsn string) (*gorm.DB, error) {
	switch dbType {
	case config.TypeSQLite:
		return OpenSQLite(dsn)
	case config.TypeMySql:
		return OpenMySql(dsn)
	case config.TypePostgreSQL:
		return OpenPostgreSQL(dsn)
	case config.TypeSqlServer:
		return OpenSqlServer(dsn)
	}
	return nil, errors.New(`不支持的数据库类型 "` + string(dbType) + `"`)
}

func OpenSQLite(filename string) (*gorm.DB, error) {
	if err := EnsureDir(filepath.Dir(filename)); err != nil {
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

func GetDsn(dbType config.DBType, host string, portStr string, dbName string, user string, password string) (string, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", errors.New("port 值有误")
	}

	switch dbType {
	case config.TypeMySql:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName), nil
	case config.TypePostgreSQL:
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbName), nil
	case config.TypeSqlServer:
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, host, port, dbName), nil
	}

	return "", errors.New(`不支持的数据库类型 "` + string(dbType) + `"`)
}
