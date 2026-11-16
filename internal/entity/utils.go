package entity

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
