package entity

func ContainsComment(comments []Comment, targetID uint) bool {
	for _, c := range comments {
		if c.ID == targetID {
			return true
		}
	}
	return false
}

func ContainsCommentDTO(comments []CommentDTO, targetID uint) bool {
	for _, c := range comments {
		if c.ID == targetID {
			return true
		}
	}
	return false
}
