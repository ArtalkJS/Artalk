package entity

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Content string

	PageKey  string `gorm:"index;size:255"`
	SiteName string `gorm:"index;size:255"`

	UserID     uint `gorm:"index"`
	IsVerified bool `gorm:"default:false"`
	UA         string
	IP         string

	Rid uint `gorm:"index"` // Parent Node ID

	IsCollapsed bool `gorm:"default:false"`
	IsPending   bool `gorm:"default:false"`
	IsPinned    bool `gorm:"default:false"`

	VoteUp   int
	VoteDown int

	RootID uint  `gorm:"index"` // Root Node ID (can be derived from `Rid`)
	Page   *Page `gorm:"foreignKey:page_key;references:key"`
	User   *User `gorm:"foreignKey:id;references:user_id"`
}

func (c Comment) IsEmpty() bool {
	return c.ID == 0
}

func (c Comment) IsAllowReply() bool {
	return !c.IsCollapsed && !c.IsPending
}
