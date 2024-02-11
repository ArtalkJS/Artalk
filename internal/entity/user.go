package entity

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string `gorm:"index;size:255"`
	Email          string `gorm:"index;size:255"`
	Link           string
	Password       string
	BadgeName      string
	BadgeColor     string
	LastIP         string
	LastUA         string
	IsAdmin        bool
	ReceiveEmail   bool `gorm:"default:true"`
	TokenValidFrom sql.NullTime

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
	u.TokenValidFrom.Scan(time.Now())
	return nil
}

func (u *User) CheckPassword(input string) bool {
	if u.ID == 0 {
		return false
	}
	password := strings.TrimSpace(u.Password)
	if password == "" {
		return false
	}

	const BcryptPrefix = "(bcrypt)"
	const MD5Prefix = "(md5)"

	switch {
	case strings.HasPrefix(password, BcryptPrefix):
		if err := bcrypt.CompareHashAndPassword([]byte(password[len(BcryptPrefix):]),
			[]byte(input)); err == nil {
			return true
		}
	case strings.HasPrefix(password, MD5Prefix):
		if strings.EqualFold(password[len(MD5Prefix):],
			fmt.Sprintf("%x", md5.Sum([]byte(input)))) {
			return true
		}
	default:
		if password == input {
			return true
		}
	}

	return false
}
