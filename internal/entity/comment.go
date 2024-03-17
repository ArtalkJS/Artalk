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

	Rid uint `gorm:"index"` // 父评论 ID

	IsCollapsed bool `gorm:"default:false"` // 折叠
	IsPending   bool `gorm:"default:false"` // 待审
	IsPinned    bool `gorm:"default:false"` // 置顶

	VoteUp   int
	VoteDown int

	RootID   uint       `gorm:"index"` // 根评论 ID
	Children []*Comment `gorm:"foreignKey:root_id;references:id"`
	Page     *Page      `gorm:"foreignKey:page_key;references:key"`
	User     *User      `gorm:"foreignKey:id;references:user_id"`
}

func (c Comment) IsEmpty() bool {
	return c.ID == 0
}

func (c Comment) IsAllowReply() bool {
	return !c.IsCollapsed && !c.IsPending
}
