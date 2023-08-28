package db

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"gorm.io/gorm"
)

func MigrateModels(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&entity.Site{}, &entity.Page{}, &entity.User{},
		&entity.Comment{}, &entity.Notify{}, &entity.Vote{}) // 注意表的创建顺序，因为有关联字段
}
