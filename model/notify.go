package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
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

	_Comment     Comment
	_CommentOnce sync.Once
}

func (n Notify) IsEmpty() bool {
	return n.ID == 0
}

func (n *Notify) FetchComment() Comment {
	if n._Comment.IsEmpty() {
		n._CommentOnce.Do(func() {
			comment := FindComment(n.CommentID)
			n._Comment = comment
		})
	}

	return n._Comment
}

func (n *Notify) SetComment(comment Comment) {
	n._Comment = comment
}

// 获取接收通知的用户
func (n *Notify) FetchUser() User {
	return FindUserByID(n.UserID)
}

func (n *Notify) GetParentComment() Comment {
	comment := n.FetchComment()
	if comment.Rid == 0 {
		return Comment{}
	}

	return FindComment(comment.Rid)
}

// 操作时的验证密钥（判断是否本人操作）
func (n *Notify) GenerateKey() {
	n.Key = lib.GetMD5Hash(fmt.Sprintf("%v %v %v", n.UserID, n.CommentID, time.Now().Unix()))
}

func (n *Notify) GetReadLink() string {
	c := n.FetchComment()

	return c.GetLinkToReply(n.Key)
}

func (n *Notify) SetInitial() error {
	n.IsRead = false
	n.IsEmailed = false
	return DB().Save(n).Error
}

func (n *Notify) SetRead() error {
	n.IsRead = true
	nowTime := time.Now()
	n.ReadAt = &nowTime
	return DB().Save(n).Error
}

func (n *Notify) SetEmailed() error {
	n.IsEmailed = true
	nowTime := time.Now()
	n.EmailAt = &nowTime
	return DB().Save(n).Error
}

type CookedNotify struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	CommentID uint   `json:"comment_id"`
	IsRead    bool   `json:"is_read"`
	IsEmailed bool   `json:"is_emailed"`
	ReadLink  string `json:"read_link"`
}

func (n *Notify) ToCooked() CookedNotify {
	return CookedNotify{
		ID:        n.ID,
		UserID:    n.UserID,
		CommentID: n.CommentID,
		IsRead:    n.IsRead,
		IsEmailed: n.IsEmailed,
		ReadLink:  n.GetReadLink(),
	}
}
