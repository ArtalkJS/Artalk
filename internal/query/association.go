package query

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

// ===============
//  Comment
// ===============

func FetchUserForComment(c *entity.Comment) entity.User {
	return FindUserByID(c.UserID)
}

func FetchPageForComment(c *entity.Comment) entity.Page {
	return FindPage(c.PageKey, c.SiteName)
}

func FetchSiteForComment(c *entity.Comment) entity.Site {
	return FindSite(c.SiteName)
}

// ===============
//  Page
// ===============

func FetchSiteForPage(p *entity.Page) entity.Site {
	return FindSite(p.SiteName)
}

// ===============
//  Notify
// ===============

func FetchCommentForNotify(n *entity.Notify) entity.Comment {
	return FindComment(n.CommentID)
}

// 获取接收通知的用户
func FetchUserForNotify(n *entity.Notify) entity.User {
	return FindUserByID(n.UserID)
}
