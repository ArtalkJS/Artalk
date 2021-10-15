package model

import (
	"strings"

	"gorm.io/gorm"
)

type VoteType string

const (
	VoteTypeCommentUp   VoteType = "comment_up"
	VoteTypeCommentDown VoteType = "comment_down"
	VoteTypePageUp      VoteType = "page_up"
	VoteTypePageDown    VoteType = "page_down"
)

type Vote struct {
	gorm.Model

	TargetID uint     `gorm:"index"` // 投票对象
	Type     VoteType `gorm:"index"`

	UserID uint `gorm:"index"` // 投票者
	UA     string
	IP     string
}

func (v *Vote) IsEmpty() bool {
	return v.ID == 0
}

func (v *Vote) IsUp() bool {
	return strings.HasSuffix(string(v.Type), "_up")
}
