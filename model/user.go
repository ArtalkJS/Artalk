package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string `gorm:"index"`
	Email      string `gorm:"index"`
	Link       string
	Password   string
	BadgeName  string
	BadgeColor string
	LastIP     string
	LastUA     string
	IsAdmin    bool

	// 配置文件中添加的
	IsInConf bool
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}

type CookedUser struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Link       string `json:"link"`
	BadgeName  string `json:"badge_name"`
	BadgeColor string `json:"badge_color"`
	IsAdmin    bool   `json:"is_admin"`
}

func (u User) ToCooked() CookedUser {
	return CookedUser{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Link:       u.Link,
		BadgeName:  u.BadgeName,
		BadgeColor: u.BadgeColor,
		IsAdmin:    u.IsAdmin,
	}
}
