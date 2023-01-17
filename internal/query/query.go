package query

import "github.com/ArtalkJS/Artalk/internal/entity"

func GetUserAllCommentIDs(userID uint) []uint {
	userAllCommentIDs := []uint{}
	DB().Model(&entity.Comment{}).Select("id").Where("user_id = ?", userID).Find(&userAllCommentIDs)
	return userAllCommentIDs
}
