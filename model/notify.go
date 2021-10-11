package model

import (
	"fmt"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"gorm.io/gorm"
)

type Notify struct {
	gorm.Model

	UserID    uint `gorm:"index"` // 通知对象
	CommentID uint `gorm:"index"` // 待查看的评论

	IsRead    bool
	ReadAt    time.Time
	IsEmailed bool
	EmailAt   time.Time

	Key string `gorm:"index"`

	Comment Comment `gorm:"foreignKey:CommentID;references:ID"`
}

func (n Notify) IsEmpty() bool {
	return n.ID == 0
}

func (n *Notify) FetchComment() Comment {
	if !n.Comment.IsEmpty() {
		return n.Comment
	}

	var comment Comment
	lib.DB.First(&comment, n.CommentID)

	n.Comment = comment
	return comment
}

func (n *Notify) GetParentComment() Comment {
	comment := n.FetchComment()
	if comment.Rid == 0 {
		return Comment{}
	}

	pComment := FindComment(comment.Rid, comment.SiteName)
	return pComment
}

// 操作时的验证密钥（判断是否本人操作）
func (n *Notify) GenerateKey() {
	n.Key = lib.GetMD5Hash(fmt.Sprintf("%v %v %v", n.UserID, n.CommentID, time.Now().Unix()))
}

func (n *Notify) GetReadLink() string {
	c := n.FetchComment()

	if !lib.ValidateURL(c.PageKey) {
		return ""
	}

	return lib.AddQueryToURL(c.PageKey, map[string]string{
		"atk_comment":    fmt.Sprintf("%d", c.ID),
		"atk_notify_key": n.Key,
	})
}

func (n *Notify) SetInitial() error {
	n.IsRead = false
	n.IsEmailed = false
	return lib.DB.Save(n).Error
}

func (n *Notify) SetRead() error {
	n.IsRead = true
	n.ReadAt = time.Now()
	return lib.DB.Save(n).Error
}

func (n *Notify) SetEmailed() error {
	n.IsEmailed = true
	n.EmailAt = time.Now()
	return lib.DB.Save(n).Error
}
