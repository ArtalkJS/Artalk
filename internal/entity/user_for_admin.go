package entity

type CookedUserForAdmin struct {
	CookedUser
	LastIP       string `json:"last_ip"`
	LastUA       string `json:"last_ua"`
	IsInConf     bool   `json:"is_in_conf"`
	CommentCount int64  `json:"comment_count"`
}
