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
