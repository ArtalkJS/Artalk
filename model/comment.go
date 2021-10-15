package model

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Content string

	PageKey  string `gorm:"index"`
	SiteName string `gorm:"index"`

	UserID uint `gorm:"index"`
	UA     string
	IP     string

	Rid uint `gorm:"index"` // 父评论 ID

	IsCollapsed bool // 折叠
	IsPending   bool // 待审

	User User `gorm:"foreignKey:UserID;references:ID"`
	Page Page `gorm:"foreignKey:PageKey;references:Key"`
	Site Site `gorm:"foreignKey:SiteName;references:Name"`

	VoteUp   int
	VoteDown int
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
	lib.DB.Where("page_key = ?", c.PageKey).First(&page)

	c.Page = page
	return page
}

func (c *Comment) FetchSite() Site {
	if !c.Site.IsEmpty() {
		return c.Site
	}

	var site Site
	lib.DB.Where("name = ?", c.SiteName).First(&site)

	c.Site = site
	return site
}

func (c Comment) FetchChildren(filters ...func(db *gorm.DB) *gorm.DB) []Comment {
	children := []Comment{}
	fetchChildrenOnce(&children, c, filters...) // TODO: children 数量限制
	return children
}

func fetchChildrenOnce(src *[]Comment, parentComment Comment, filters ...func(db *gorm.DB) *gorm.DB) {
	children := []Comment{}
	lib.DB.Scopes(filters...).Where("rid = ?", parentComment.ID).Order("created_at ASC").Find(&children)

	for _, child := range children {
		*src = append(*src, child)
		fetchChildrenOnce(src, child, filters...) // loop
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
	BadgeName      string `json:"badge_name"`
	BadgeColor     string `json:"badge_color"`
	Visible        bool   `json:"visible"`
	VoteUp         int    `json:"vote_up"`
	VoteDown       int    `json:"vote_down"`
	PageKey        string `json:"page_key"`
	SiteName       string `json:"site_name"`
}

func (c *Comment) ToCooked() CookedComment {
	user := c.FetchUser()
	//page := c.FetchPage()

	return CookedComment{
		ID:             c.ID,
		Content:        c.Content,
		Nick:           user.Name,
		EmailEncrypted: lib.GetMD5Hash(user.Email),
		Link:           user.Link,
		UA:             c.UA,
		Date:           c.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		IsCollapsed:    c.IsCollapsed,
		IsPending:      c.IsPending,
		IsAllowReply:   c.IsAllowReply(),
		Rid:            c.Rid,
		BadgeName:      user.BadgeName,
		BadgeColor:     user.BadgeColor,
		Visible:        true,
		VoteUp:         c.VoteUp,
		VoteDown:       c.VoteDown,
		PageKey:        c.PageKey,
		SiteName:       c.SiteName,
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
	BadgeName      string `json:"badge_name"`
	BadgeColor     string `json:"badge_color"`
	PageKey        string `json:"page_key"`
	SiteName       string `json:"site_name"`
}

func (c *Comment) ToCookedForEmail() CookedCommentForEmail {
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
		BadgeName:      user.BadgeName,
		BadgeColor:     user.BadgeColor,
		PageKey:        c.PageKey,
		SiteName:       c.SiteName,
	}
}

func (c *Comment) SpamCheck(echoCtx echo.Context) {
	setPending := func() {
		if c.IsPending {
			return
		}

		// 改为待审状态
		c.IsPending = true
		lib.DB.Save(c)
	}

	siteURL := ""
	if c.SiteName != "" {
		site := FindSite(c.SiteName)
		siteURL = site.ToCooked().FirstUrl
	}
	if siteURL == "" { // 从 referer 中提取网站
		if pr, err := url.Parse(echoCtx.Request().Referer()); err == nil && pr.Scheme != "" && pr.Host != "" {
			siteURL = fmt.Sprintf("%s://%s", pr.Scheme, pr.Host)
		}
	}

	user := c.FetchUser()

	// akismet
	akismetKey := strings.TrimSpace(config.Instance.Moderator.AkismetKey)
	if akismetKey != "" {
		isOK, err := lib.AntiSpamCheck_Akismet(&lib.AkismetParams{
			Blog:               siteURL,
			UserIP:             echoCtx.RealIP(),
			UserAgent:          echoCtx.Request().UserAgent(),
			CommentType:        "comment",
			CommentAuthor:      user.Name,
			CommentAuthorEmail: user.Email,
			CommentContent:     c.Content,
		}, akismetKey)
		if err != nil {
			logrus.Error("akismet 垃圾检测错误 ", err)
		}
		if !isOK {
			setPending()
		}
	}

	// TODO:关键字过滤
}
