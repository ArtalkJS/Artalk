package dao

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

func (dao *Dao) MigrateModels() {
	// Migrate the schema
	dao.DB().AutoMigrate(&entity.Site{}, &entity.Page{}, &entity.User{},
		&entity.Comment{}, &entity.Notify{}, &entity.Vote{}) // 注意表的创建顺序，因为有关联字段
}
