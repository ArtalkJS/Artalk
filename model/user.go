package model

import (
	"gorm.io/gorm"
)

type UserType string

const (
	UserBanned UserType = "banned"
	UserAdmin  UserType = "admin"
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
	Type       UserType `gorm:"index"`
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}

func (u User) IsAdmin() bool {
	return u.Type == UserAdmin
}

type CookedUser struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Link       string   `json:"link"`
	BadgeName  string   `json:"badge_name"`
	BadgeColor string   `json:"badge_color"`
	IsAdmin    bool     `json:"is_admin"`
	Type       UserType `json:"type"`
}

func (u User) ToCooked() CookedUser {
	return CookedUser{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Link:       u.Link,
		BadgeName:  u.BadgeName,
		BadgeColor: u.BadgeColor,
		IsAdmin:    u.IsAdmin(),
		Type:       u.Type,
	}
}
