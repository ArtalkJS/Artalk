package model

import (
	"fmt"
	"sync"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Content string

	PageKey  string `gorm:"index;size:255"`
	SiteName string `gorm:"index;size:255"`

	UserID uint `gorm:"index"`
	UA     string
	IP     string

	Rid uint `gorm:"index"` // 父评论 ID

	IsCollapsed bool `gorm:"default:false"` // 折叠
	IsPending   bool `gorm:"default:false"` // 待审
	IsPinned    bool `gorm:"default:false"` // 置顶

	_User User
	_Page Page
	_Site Site

	_UserOnce sync.Once
	_PageOnce sync.Once
	_SiteOnce sync.Once

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
	if c._User.IsEmpty() {
		c._UserOnce.Do(func() {
			user := FindUserByID(c.UserID)
			c._User = user
		})
	}

	return c._User
}

func (c *Comment) FetchPage() Page {
	if c._Page.IsEmpty() {
		c._PageOnce.Do(func() {
			page := FindPage(c.PageKey, c.SiteName)
			c._Page = page
		})
	}

	return c._Page
}

func (c *Comment) FetchSite() Site {
	if c._Site.IsEmpty() {
		c._SiteOnce.Do(func() {
			site := FindSite(c.SiteName)
			c._Site = site
		})
	}

	return c._Site
}

// 获取评论回复链接
func (c *Comment) GetLinkToReply(notifyKey ...string) string {
	page := c.FetchPage()
	rawURL := page.GetAccessibleURL()

	// 请求 query
	queryMap := map[string]string{
		"atk_comment": fmt.Sprintf("%d", c.ID),
	}

	// atk_notify_key
	if len(notifyKey) > 0 {
		queryMap["atk_notify_key"] = notifyKey[0]
	}

	return lib.AddQueryToURL(rawURL, queryMap)
}

type CookedComment struct {
	ID             uint   `json:"id"`
	Content        string `json:"content"`
	ContentMarked  string `json:"content_marked"`
	UserID         uint   `json:"user_id"`
	Nick           string `json:"nick"`
	EmailEncrypted string `json:"email_encrypted"`
	Link           string `json:"link"`
	UA             string `json:"ua"`
	Date           string `json:"date"`
	IsCollapsed    bool   `json:"is_collapsed"`
	IsPending      bool   `json:"is_pending"`
	IsPinned       bool   `json:"is_pinned"`
	IsAllowReply   bool   `json:"is_allow_reply"`
	Rid            uint   `json:"rid"`
	BadgeName      string `json:"badge_name"`
	BadgeColor     string `json:"badge_color"`
	Visible        bool   `json:"visible"`
	VoteUp         int    `json:"vote_up"`
	VoteDown       int    `json:"vote_down"`
	PageKey        string `json:"page_key"`
	PageURL        string `json:"page_url"`
	SiteName       string `json:"site_name"`
}

func (c *Comment) ToCooked() CookedComment {
	user := c.FetchUser()
	page := c.FetchPage()

	markedContent, _ := lib.Marked(c.Content)

	return CookedComment{
		ID:             c.ID,
		Content:        c.Content,
		ContentMarked:  markedContent,
		UserID:         c.UserID,
		Nick:           user.Name,
		EmailEncrypted: lib.GetMD5Hash(user.Email),
		Link:           user.Link,
		UA:             c.UA,
		Date:           c.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		IsCollapsed:    c.IsCollapsed,
		IsPending:      c.IsPending,
		IsPinned:       c.IsPinned,
		IsAllowReply:   c.IsAllowReply(),
		Rid:            c.Rid,
		BadgeName:      user.BadgeName,
		BadgeColor:     user.BadgeColor,
		Visible:        true,
		VoteUp:         c.VoteUp,
		VoteDown:       c.VoteDown,
		PageKey:        c.PageKey,
		PageURL:        page.GetAccessibleURL(),
		SiteName:       c.SiteName,
	}
}

type CookedCommentForEmail struct {
	CookedComment
	Content    string     `json:"content"`
	ContentRaw string     `json:"content_raw"`
	Nick       string     `json:"nick"`
	Email      string     `json:"email"`
	IP         string     `json:"ip"`
	Datetime   string     `json:"datetime"`
	Date       string     `json:"date"`
	Time       string     `json:"time"`
	PageKey    string     `json:"page_key"`
	PageTitle  string     `json:"page_title"`
	Page       CookedPage `json:"page"`
	SiteName   string     `json:"site_name"`
	Site       CookedSite `json:"site"`
}

func (c *Comment) ToCookedForEmail() CookedCommentForEmail {
	user := c.FetchUser()
	page := c.FetchPage()
	site := c.FetchSite()
	content, _ := lib.Marked(c.Content)

	return CookedCommentForEmail{
		Content:    content,
		ContentRaw: c.Content,
		Nick:       user.Name,
		Email:      user.Email,
		IP:         c.IP,
		Datetime:   c.CreatedAt.Local().Format("2006-01-02 15:04:05"),
		Date:       c.CreatedAt.Local().Format("2006-01-02"),
		Time:       c.CreatedAt.Local().Format("15:04:05"),
		PageKey:    c.PageKey,
		PageTitle:  page.Title,
		Page:       page.ToCooked(),
		SiteName:   c.SiteName,
		Site:       site.ToCooked(),
		CookedComment: CookedComment{
			ID:             c.ID,
			EmailEncrypted: lib.GetMD5Hash(user.Email),
			Link:           user.Link,
			UA:             c.UA,
			IsCollapsed:    c.IsCollapsed,
			IsPending:      c.IsPending,
			IsPinned:       c.IsPinned,
			IsAllowReply:   c.IsAllowReply(),
			Rid:            c.Rid,
			BadgeName:      user.BadgeName,
			BadgeColor:     user.BadgeColor,
		},
	}
}

func (c *Comment) ToArtran() Artran {
	user := c.FetchUser()
	page := c.FetchPage()
	site := c.FetchSite()

	return Artran{
		ID:            lib.ToString(c.ID),
		Rid:           lib.ToString(c.Rid),
		Content:       c.Content,
		UA:            c.UA,
		IP:            c.IP,
		IsCollapsed:   lib.ToString(c.IsCollapsed),
		IsPending:     lib.ToString(c.IsPending),
		IsPinned:      lib.ToString(c.IsPinned),
		VoteUp:        lib.ToString(c.VoteUp),
		VoteDown:      lib.ToString(c.VoteDown),
		CreatedAt:     c.CreatedAt.String(),
		UpdatedAt:     c.UpdatedAt.String(),
		Nick:          user.Name,
		Email:         user.Email,
		Link:          user.Link,
		BadgeName:     user.BadgeName,
		BadgeColor:    user.BadgeColor,
		PageKey:       page.Key,
		PageTitle:     page.Title,
		PageAdminOnly: lib.ToString(page.AdminOnly),
		SiteName:      site.Name,
		SiteUrls:      site.Urls,
	}
}
