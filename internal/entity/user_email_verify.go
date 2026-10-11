package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserEmailVerify struct {
	gorm.Model
	Email string `gorm:"index;size:255"`
	Code  string
	// Persisted per challenge; zero preserves unused legacy codes on upgrade.
	Attempts  int `gorm:"not null;default:0"`
	ExpiresAt time.Time
	IP        string
	UA        string
}
