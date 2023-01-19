package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/ArtalkJS/Artalk/internal/utils"
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

	Comment      Comment   `gorm:"-"`
	Once_Comment sync.Once `gorm:"-"`
}

func (n Notify) IsEmpty() bool {
	return n.ID == 0
}

func (n *Notify) SetComment(comment Comment) {
	n.Comment = comment
}

// 操作时的验证密钥（判断是否本人操作）
func (n *Notify) GenerateKey() {
	n.Key = utils.GetMD5Hash(fmt.Sprintf("%v %v %v", n.UserID, n.CommentID, time.Now().Unix()))
}
