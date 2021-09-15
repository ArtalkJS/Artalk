package model

import "gorm.io/gorm"

type PageType string

const ()

type Page struct {
	gorm.Model
	Key       string `gorm:"uniqueIndex"`
	AdminOnly bool
	Type      PageType
}

func (p Page) IsEmpty() bool {
	return p.ID == 0
}

type CookedPage struct {
	ID        uint   `json:"id"`
	AdminOnly bool   `json:"admin_only"`
	PageKey   string `json:"page_key"`
}

func (p Page) ToCooked() CookedPage {
	return CookedPage{
		ID:        p.ID,
		AdminOnly: p.AdminOnly,
		PageKey:   p.Key,
	}
}
