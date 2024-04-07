// This file contains the exposed methods for `comments_get` package
package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

type FindOptions struct {
	Offset int
	Limit  int
	Nested bool
}

// Find comments by options
func FindComments(dao *dao.Dao, opts QueryOptions, pg FindOptions) ([]entity.CookedComment, int64, int64) {
	// Shared scopes
	// Generated where conditions by options
	var scopes []func(*gorm.DB) *gorm.DB
	scopes = append(scopes, ConvertGormScopes(GetQueryScopes(dao, opts))...)
	scopes = append(scopes, func(d *gorm.DB) *gorm.DB {
		return d.Preload("User").Preload("Page").Preload("Page.Site")
	})

	// First query
	var comments []*entity.Comment
	dao.DB().Model(&entity.Comment{}).
		Scopes(scopes...).
		Scopes(func(d *gorm.DB) *gorm.DB {
			if pg.Nested {
				d.Scopes(OnlyRoot()) // Nested mode get only the root comments
			}
			return d
		}).
		Order(GetSortSQL(opts.Scope, opts.SortBy)).
		Offset(pg.Offset).
		Limit(pg.Limit).
		Find(&comments)

	// Subsequent query
	cooked := dao.CookAllComments(comments)
	if pg.Nested {
		cooked = findNestedChildren(dao, cooked, scopes)
	} else {
		cooked = findFlatLinkedComments(dao, cooked, scopes)
	}

	// Get count
	var count int64
	var rootsCount int64
	{
		dao.DB().Model(&entity.Comment{}).Scopes(scopes...).Count(&count) // Note: Count will omit preloads
		dao.DB().Model(&entity.Comment{}).Scopes(scopes...).Scopes(OnlyRoot()).Count(&rootsCount)
	}

	return cooked, count, rootsCount
}
