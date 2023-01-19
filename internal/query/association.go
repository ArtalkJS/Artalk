package query

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

// ===============
//  Comment
// ===============

func FetchUserForComment(c *entity.Comment) entity.User {
	if c.User.IsEmpty() {
		c.Once_User.Do(func() {
			user := FindUserByID(c.UserID)
			c.User = user
		})
	}

	return c.User
}

func FetchPageForComment(c *entity.Comment) entity.Page {
	if c.Page.IsEmpty() {
		c.Once_Page.Do(func() {
			page := FindPage(c.PageKey, c.SiteName)
			c.Page = page
		})
	}

	return c.Page
}

func FetchSiteForComment(c *entity.Comment) entity.Site {
	if c.Site.IsEmpty() {
		c.Once_Site.Do(func() {
			site := FindSite(c.SiteName)
			c.Site = site
		})
	}

	return c.Site
}

// ===============
//  Page
// ===============

func FetchSiteForPage(p *entity.Page) entity.Site {
	if p.Site.IsEmpty() {
		p.Once_Site.Do(func() {
			site := FindSite(p.SiteName)
			p.Site = site
		})
	}

	return p.Site
}

// ===============
//  Notify
// ===============

func FetchCommentForNotify(n *entity.Notify) entity.Comment {
	if n.Comment.IsEmpty() {
		n.Once_Comment.Do(func() {
			comment := FindComment(n.CommentID)
			n.Comment = comment
		})
	}

	return n.Comment
}

// 获取接收通知的用户
func FetchUserForNotify(n *entity.Notify) entity.User {
	return FindUserByID(n.UserID)
}
