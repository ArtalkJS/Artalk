package model

import "gorm.io/gorm"

type PageType string

const (
	PageCommentClosed PageType = "comment_closed"
	PageOnlyAdmin     PageType = "only_admin"
)

type Page struct {
	gorm.Model
	Key  string `gorm:"uniqueIndex"`
	Type PageType
}

func (p Page) IsEmpty() bool {
	return p.ID == 0
}

type CookedPage struct {
	ID              uint   `json:"id"`
	IsCommentClosed bool   `json:"is_comment_closed"`
	PageKey         string `json:"page_key"`
}

func (p Page) ToCooked() CookedPage {
	return CookedPage{
		ID:              p.ID,
		IsCommentClosed: p.Type == PageCommentClosed,
		PageKey:         p.Key,
	}
}
