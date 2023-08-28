package dao

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

func (dao *Dao) FindCreateSite(siteName string) entity.Site {
	site := dao.FindSite(siteName)
	if site.IsEmpty() {
		site = dao.NewSite(siteName, "")
	}
	return site
}

func (dao *Dao) FindCreatePage(pageKey string, pageTitle string, siteName string) entity.Page {
	page := dao.FindPage(pageKey, siteName)
	if page.IsEmpty() {
		page = dao.NewPage(pageKey, pageTitle, siteName)
	}
	return page
}

func (dao *Dao) FindCreateUser(name string, email string, link string) entity.User {
	user := dao.FindUser(name, email)
	if user.IsEmpty() {
		user = dao.NewUser(name, email, link) // save a new user
	}
	return user
}

func (dao *Dao) FindCreateNotify(userID uint, lookCommentID uint) entity.Notify {
	notify := dao.FindNotify(userID, lookCommentID)
	if notify.IsEmpty() {
		notify = dao.NewNotify(userID, lookCommentID)
	}
	return notify
}
