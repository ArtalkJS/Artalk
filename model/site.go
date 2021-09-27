package model

import (
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
	Url  string
}

func (s Site) IsEmpty() bool {
	return s.ID == 0
}

type CookedSite struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (s Site) ToCooked() CookedSite {
	return CookedSite{
		ID:   s.ID,
		Name: s.Name,
		Url:  s.Url,
	}
}
