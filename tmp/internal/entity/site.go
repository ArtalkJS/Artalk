package entity

import (
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;size:255"`
	Urls string
}

func (s Site) IsEmpty() bool {
	return s.ID == 0
}
