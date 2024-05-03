package entity

import (
	"time"

	"gorm.io/gorm"
)

type AuthIdentity struct {
	gorm.Model
	Provider    string // local, email, oauth, github
	RemoteUID   string `gorm:"column:remote_uid;index;size:255"`
	UserID      uint   `gorm:"index"`
	Token       string
	ConfirmedAt *time.Time
	ExpiresAt   *time.Time
}

func (n AuthIdentity) IsEmpty() bool {
	return n.ID == 0
}
