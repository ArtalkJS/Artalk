package model

import "gorm.io/gorm"

type PageType string

const (
	PageClosed PageType = "closed"
)

type Page struct {
	gorm.Model
	Key  string
	Type PageType
}

func (p Page) IsEmpty() bool {
	return p.ID == 0
}
