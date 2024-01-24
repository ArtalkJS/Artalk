package entity

import (
	"golang.org/x/crypto/bcrypt"
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
	ReceiveEmail bool `gorm:"default:true"`

	// 配置文件中添加的
	IsInConf bool
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}

func (u *User) SetPasswordEncrypt(password string) (err error) {
	var encrypted []byte
	if encrypted, err = bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	); err != nil {
		return err
	}
	u.Password = "(bcrypt)" + string(encrypted)
	return nil
}
