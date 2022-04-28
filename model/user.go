package model

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"index;size:255"`
	Email        string `gorm:"index;size:255"`
	Link         string
	Password     string
	BadgeName    string
	BadgeColor   string
	LastIP       string
	LastUA       string
	IsAdmin      bool
	SiteNames    string
	ReceiveEmail bool `gorm:"default:true"`

	// 配置文件中添加的
	IsInConf bool
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}

type CookedUser struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Link         string   `json:"link"`
	BadgeName    string   `json:"badge_name"`
	BadgeColor   string   `json:"badge_color"`
	IsAdmin      bool     `json:"is_admin"`
	SiteNames    []string `json:"site_names"`
	SiteNamesRaw string   `json:"site_names_raw"`
	ReceiveEmail bool     `json:"receive_email"`
}

func (u User) ToCooked() CookedUser {
	splitSites := lib.RemoveBlankStrings(lib.SplitAndTrimSpace(u.SiteNames, ","))

	return CookedUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Link:         u.Link,
		BadgeName:    u.BadgeName,
		BadgeColor:   u.BadgeColor,
		IsAdmin:      u.IsAdmin,
		SiteNames:    splitSites,
		SiteNamesRaw: u.SiteNames,
		ReceiveEmail: u.ReceiveEmail,
	}
}
