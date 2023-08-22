package entity

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Notify struct {
	gorm.Model

	UserID    uint `gorm:"index"` // 通知对象 (接收通知的用户 ID)
	CommentID uint `gorm:"index"` // 待查看的评论

	IsRead    bool
	ReadAt    *time.Time
	IsEmailed bool
	EmailAt   *time.Time

	Key string `gorm:"index;size:255"`
}

func (n Notify) IsEmpty() bool {
	return n.ID == 0
}

// 操作时的验证密钥（判断是否本人操作）
func (n *Notify) GenerateKey() {
	letterRunes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	n.Key = string(b)
}
