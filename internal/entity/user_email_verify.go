package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserEmailVerify struct {
	gorm.Model
	Email     string `gorm:"index;size:255"`
	Code      string
	ExpiresAt time.Time
	IP        string
	UA        string
}
