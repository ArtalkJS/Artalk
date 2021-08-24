package model

import "gorm.io/gorm"

type CommentType string

const (
	CommentCollapsed CommentType = "collapsed"
	CommentPending   CommentType = "pengding"
	CommentDeleted   CommentType = "deleted"
)

type Comment struct {
	gorm.Model
	Content string

	UserID uint
	PageID uint
	User   User `gorm:"foreignKey:UserID"`
	Page   Page `gorm:"foreignKey:PageID"`

	Rid  uint
	UA   string
	IP   string
	Type CommentType
}
