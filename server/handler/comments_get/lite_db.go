package comments_get

import "gorm.io/gorm"

// LiteDB is a simplified form of gorm.DB, made to make handling database tasks easier.
// LiteDB only has `Where` and `Scopes` methods, which simplifies the creation of queries.
// This simpler structure is especially helpful when writing test cases,
// as it can make the code more straightforward and easier to comprehend.
type liteDB interface {
	Where(query interface{}, args ...interface{}) liteDB
	Scopes(funcs ...func(liteDB) liteDB) liteDB
}

func ConvertGormScopes(funcs ...func(liteDB) liteDB) []func(*gorm.DB) *gorm.DB {
	var ret []func(*gorm.DB) *gorm.DB
	for _, fn := range funcs {
		ret = append(ret, func(db *gorm.DB) *gorm.DB {
			return fn(&liteDbImpl{db: db}).(*liteDbImpl).db
		})
	}
	return ret
}

type liteDbImpl struct {
	db *gorm.DB
}

func (l *liteDbImpl) Where(query interface{}, args ...interface{}) liteDB {
	l.db = l.db.Where(query, args...)
	return l
}

func (l *liteDbImpl) Scopes(funcs ...func(liteDB) liteDB) liteDB {
	for _, fn := range funcs {
		l = fn(l).(*liteDbImpl)
	}
	return l
}
