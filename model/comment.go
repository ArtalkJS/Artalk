package model

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"gorm.io/gorm"
)

type CommentType string

type Comment struct {
	gorm.Model
	Content string

	UserID  uint   `gorm:"index"`
	PageKey string `gorm:"index"`
	User    User   `gorm:"foreignKey:UserID;references:ID"`
	Page    Page   `gorm:"foreignKey:PageKey;references:Key"`

	Rid uint `gorm:"index"`
	UA  string
	IP  string

	IsCollapsed bool
	IsPending   bool
}

func (c Comment) IsEmpty() bool {
	return c.ID == 0
}

func (c Comment) IsAllowReply() bool {
	return !c.IsCollapsed && !c.IsPending
}

func (c *Comment) FetchUser() User {
	if !c.User.IsEmpty() {
		return c.User
	}

	// TODO: 先从 Redis 查询
	var user User
	lib.DB.First(&user, c.UserID)

	c.User = user
	return user
}

func (c *Comment) FetchPage() Page {
	if !c.Page.IsEmpty() {
		return c.Page
	}

	var page Page
	lib.DB.Where(&Page{Key: c.PageKey}).First(&page)

	c.Page = page
	return page
}

func (c Comment) FetchChildren(filter func(db *gorm.DB) *gorm.DB) []Comment {
	children := []Comment{}
	fetchChildrenOnce(&children, c, filter) // TODO: children 数量限制
	return children
}

func fetchChildrenOnce(src *[]Comment, parentComment Comment, filter func(db *gorm.DB) *gorm.DB) {
	children := []Comment{}
	lib.DB.Scopes(filter).Where("rid = ?", parentComment.ID).Order("created_at ASC").Find(&children)

	for _, child := range children {
		*src = append(*src, child)
		fetchChildrenOnce(src, child, filter) // loop
	}
}

type CookedComment struct {
	ID             uint   `json:"id"`
	Content        string `json:"content"`
	Nick           string `json:"nick"`
	EmailEncrypted string `json:"email_encrypted"`
	Link           string `json:"link"`
	UA             string `json:"ua"`
	Date           string `json:"date"`
	IsCollapsed    bool   `json:"is_collapsed"`
	IsPending      bool   `json:"is_pending"`
	IsAllowReply   bool   `json:"is_allow_reply"`
	Rid            uint   `json:"rid"`
}

func (c Comment) ToCooked() CookedComment {
	user := c.FetchUser()
	//page := c.FetchPage()

	return CookedComment{
		ID:             c.ID,
		Content:        c.Content,
		Nick:           user.Name,
		EmailEncrypted: lib.GetMD5Hash(user.Email),
		Link:           user.Link,
		UA:             c.UA,
		Date:           c.CreatedAt.Local().String(),
		IsCollapsed:    c.IsCollapsed,
		IsPending:      c.IsPending,
		IsAllowReply:   c.IsAllowReply(),
		Rid:            c.Rid,
	}
}

type CookedCommentForEmail struct {
	ID             uint   `json:"id"`
	ContentRaw     string `json:"content_raw"`
	Content        string `json:"content"`
	Nick           string `json:"nick"`
	Email          string `json:"email"`
	EmailEncrypted string `json:"email_encrypted"`
	Link           string `json:"link"`
	UA             string `json:"ua"`
	Datetime       string `json:"datetime"`
	Date           string `json:"date"`
	Time           string `json:"time"`
	IsCollapsed    bool   `json:"is_collapsed"`
	IsPending      bool   `json:"is_pending"`
	IsAllowReply   bool   `json:"is_allow_reply"`
	Rid            uint   `json:"rid"`
	PageKey        string `json:"page_key"`
}

func (c Comment) ToCookedForEmail() CookedCommentForEmail {
	user := c.FetchUser()
	content, _ := lib.Marked(c.Content)

	return CookedCommentForEmail{
		ID:             c.ID,
		ContentRaw:     c.Content,
		Content:        content,
		Nick:           user.Name,
		Email:          user.Email,
		EmailEncrypted: lib.GetMD5Hash(user.Email),
		Link:           user.Link,
		UA:             c.UA,
		Datetime:       c.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		Date:           c.CreatedAt.Local().Format("2006-01-02"),
		Time:           c.CreatedAt.Local().Format("15:04:05"),
		IsCollapsed:    c.IsCollapsed,
		IsPending:      c.IsPending,
		IsAllowReply:   c.IsAllowReply(),
		Rid:            c.Rid,
		PageKey:        c.PageKey,
	}
}
