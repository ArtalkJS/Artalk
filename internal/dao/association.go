package dao

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
)

// ===============
//  Comment
// ===============

func (dao *Dao) FetchUserForComment(c *entity.Comment) entity.User {
	return dao.FindUserByID(c.UserID)
}

func (dao *Dao) FetchPageForComment(c *entity.Comment) entity.Page {
	return dao.FindPage(c.PageKey, c.SiteName)
}

func (dao *Dao) FetchSiteForComment(c *entity.Comment) entity.Site {
	return dao.FindSite(c.SiteName)
}

// ===============
//  Page
// ===============

func (dao *Dao) FetchSiteForPage(p *entity.Page) entity.Site {
	return dao.FindSite(p.SiteName)
}

// ===============
//  Notify
// ===============

func (dao *Dao) FetchCommentForNotify(n *entity.Notify) entity.Comment {
	return dao.FindComment(n.CommentID)
}

// 获取接收通知的用户
func (dao *Dao) FetchUserForNotify(n *entity.Notify) entity.User {
	return dao.FindUserByID(n.UserID)
}
