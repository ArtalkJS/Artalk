package model

import (
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	Name   string
	Url    string
	UserID uint
}

func (s Site) IsEmpty() bool {
	return s.ID == 0
}

type CookedSite struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	UserID uint   `json:"user_id"`
}

func (s Site) ToCooked() CookedSite {
	return CookedSite{
		ID:     s.ID,
		Name:   s.Name,
		Url:    s.Url,
		UserID: s.UserID,
	}
}
