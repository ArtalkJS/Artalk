package entity

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model

	// TODO
	// For historical reasons, the naming of this field does not follow best practices.
	// Keywords should be avoided in naming, and "key" should be changed to a different name.
	// Additionally, because the unique constraint was not initially set for this field,
	// some older user databases may have duplicate "pageKey" entries due to various reasons.
	// Therefore, adding a UNIQUE constraint DIRECTLY could cause issues.
	Key string `gorm:"index;size:255"` // Page key

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
