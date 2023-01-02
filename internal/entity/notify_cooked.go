package entity

type CookedNotify struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	CommentID uint   `json:"comment_id"`
	IsRead    bool   `json:"is_read"`
	IsEmailed bool   `json:"is_emailed"`
	ReadLink  string `json:"read_link"`
}
