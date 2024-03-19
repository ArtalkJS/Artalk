package entity

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Key       string `gorm:"index;size:255"` // 页面 Key（一般为不含 hash/query 的完整 url）
	Title     string
	AdminOnly bool

	SiteName string `gorm:"index;size:255"`

	AccessibleURL string `gorm:"-"`

	VoteUp   int
	VoteDown int

	PV int

	Site *Site `gorm:"foreignKey:site_name;references:name"`
}

func (p Page) IsEmpty() bool {
	return p.ID == 0
}
