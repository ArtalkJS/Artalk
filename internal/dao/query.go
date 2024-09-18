package dao

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
	"gorm.io/gorm"
)

func (dao *Dao) GetUserAllCommentIDs(userID uint) []uint {
	userAllCommentIDs := []uint{}
	dao.DB().Model(&entity.Comment{}).Select("id").Where("user_id = ?", userID).Find(&userAllCommentIDs)
	return userAllCommentIDs
}

// Get the table name of the entity
func (dao *Dao) GetTableName(entity any) string {
	// @see https://github.com/go-gorm/gorm/issues/3603#issuecomment-709883403
	stmt := &gorm.Statement{DB: dao.DB()}
	stmt.Parse(entity)
	return stmt.Schema.Table
}
