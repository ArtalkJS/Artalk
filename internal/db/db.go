package db

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/config"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDB(conf config.DBConf) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: NewGormLogger(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.TablePrefix,
		},
	}

	var dsn string
	if conf.Dsn != "" {
		dsn = conf.Dsn
	} else {
		dsn = getDsnByConf(conf)
	}

	switch conf.Type {
	case config.TypeSQLite:
		return OpenSQLite(dsn, gormConfig)
	case config.TypeMySql:
		return OpenMySql(dsn, gormConfig)
	case config.TypePostgreSQL:
		return OpenPostgreSQL(dsn, gormConfig)
	case config.TypeMSSQL:
		return OpenSqlServer(dsn, gormConfig)
	}

	return nil, fmt.Errorf("unsupported database type %s", conf.Type)
}
