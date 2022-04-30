package model

func CookAllComments(comments []Comment) []CookedComment {
	cookedComments := []CookedComment{}
	for _, c := range comments {
		cookedComments = append(cookedComments, c.ToCooked())
	}
	return cookedComments
}

func GetUserAllCommentIDs(userID uint) []uint {
	userAllCommentIDs := []uint{}
	DB().Model(&Comment{}).Select("id").Where("user_id = ?", userID).Find(&userAllCommentIDs)
	return userAllCommentIDs
}

func ContainsComment(comments []Comment, targetID uint) bool {
	for _, c := range comments {
		if c.ID == targetID {
			return true
		}
	}
	return false
}

func ContainsCookedComment(comments []CookedComment, targetID uint) bool {
	for _, c := range comments {
		if c.ID == targetID {
			return true
		}
	}
	return false
}
