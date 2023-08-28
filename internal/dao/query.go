package dao

import "github.com/ArtalkJS/Artalk/internal/entity"

func (dao *Dao) GetUserAllCommentIDs(userID uint) []uint {
	userAllCommentIDs := []uint{}
	dao.DB().Model(&entity.Comment{}).Select("id").Where("user_id = ?", userID).Find(&userAllCommentIDs)
	return userAllCommentIDs
}
