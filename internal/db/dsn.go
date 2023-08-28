package db

import (
	"fmt"

	"github.com/ArtalkJS/Artalk/internal/config"
)

func getDsnByConf(conf config.DBConf) string {
	var dsn string

	switch conf.Type {
	case config.TypeSQLite:
		dsn = conf.File

	case config.TypePostgreSQL:
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			conf.Host,
			conf.User,
			conf.Password,
			conf.Name,
			conf.Port)

	case config.TypeMySql, config.TypeMSSQL:
		dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name,
			conf.Charset,
		)
	}

	return dsn
}
