package model

import (
	"fmt"

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
	if !c._User.IsEmpty() {
		return c._User
	}

	user := FindUserByID(c.UserID)

	c._User = user
	return user
}

func (c *Comment) FetchPage() Page {
	if !c._Page.IsEmpty() {
		return c._Page
	}

	page := FindPage(c.PageKey, c.SiteName)

	c._Page = page
	return page
}

func (c *Comment) FetchSite() Site {
	if !c._Site.IsEmpty() {
		return c._Site
	}

	site := FindSite(c.SiteName)

	c._Site = site
	return site
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

	return CookedComment{
		ID:             c.ID,
		Content:        c.Content,
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

func (c CookedComment) FetchChildrenWithCheckers(checkers ...func(*Comment) bool) []CookedComment {
	children := []CookedComment{}
	fetchChildrenOnceWithCheckers(&children, c, checkers...) // TODO: children 数量限制
	return children
}

func fetchChildrenOnceWithCheckers(src *[]CookedComment, parentComment CookedComment, checkers ...func(*Comment) bool) {
	// TODO 子评论排序问题
	children := FindCommentChildren(parentComment.ID, checkers...)

	for _, child := range children {
		*src = append(*src, child.ToCooked())
		fetchChildrenOnceWithCheckers(src, child.ToCooked(), checkers...) // loop
	}
}

// TODO 已弃用 (原因：不容易做缓存)
// func (c CookedComment) _Fetch_Children(filters ...func(db *gorm.DB) *gorm.DB) []CookedComment {
// 	children := []CookedComment{}
// 	_fetch_ChildrenOnce(&children, c, filters...) // TODO: children 数量限制
// 	return children
// }

// func _fetch_ChildrenOnce(src *[]CookedComment, parentComment CookedComment, filters ...func(db *gorm.DB) *gorm.DB) {
// 	children := []Comment{}
// 	DB().Scopes(filters...).Where("rid = ?", parentComment.ID).Order("created_at ASC").Find(&children)

// 	for _, child := range children {
// 		*src = append(*src, child.ToCooked())
// 		_fetch_ChildrenOnce(src, child.ToCooked(), filters...) // loop
// 	}
// }

type CookedCommentForEmail struct {
	CookedComment
	Content    string     `json:"content"`
	ContentRaw string     `json:"content_raw"`
	Nick       string     `json:"nick"`
	Email      string     `json:"email"`
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
