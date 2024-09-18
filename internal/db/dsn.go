package db

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/samber/lo"
)

func getDsnByConf(conf config.DBConf) string {
	var dsn string

	switch conf.Type {
	case config.TypeSQLite:
		dsn = conf.File

	case config.TypePostgreSQL:
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			conf.Host,
			conf.User,
			conf.Password,
			conf.Name,
			conf.Port,
			lo.If(conf.SSL, "require").Else("disable"),
		)

	case config.TypeMySql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&tls=%s",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name,
			conf.Charset,
			lo.If(conf.SSL, "true").Else("false"),
		)

	case config.TypeMSSQL:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name,
		)
	}

	return dsn
}
