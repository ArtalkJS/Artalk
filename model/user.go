package model

import "gorm.io/gorm"

type UserType string

const (
	UserBanned UserType = "banned"
	UserAdmin  UserType = "admin"
)

type User struct {
	gorm.Model
	Name       string `mapstructure:"name"`
	Email      string `mapstructure:"email"`
	Link       string `mapstructure:"link"`
	Password   string `mapstructure:"password"`
	BadgeName  string `mapstructure:"badge_name"`
	BadgeColor string `mapstructure:"badge_color"`
	LastIP     string
	LastUA     string
	Type       UserType
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}
