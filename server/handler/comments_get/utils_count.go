package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

// Count comments
func CountComments(dao *dao.Dao, opts QueryOptions, scopes ...func(db *gorm.DB) *gorm.DB) int64 {
	var count int64

	dao.DB().Model(&entity.Comment{}).
		Scopes(GetQueryScopes(dao, opts)).
		Scopes(scopes...).
		Count(&count)

	return count
}
