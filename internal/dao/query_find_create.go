package dao

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/utils"
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

func (dao *Dao) FindCreateUser(name string, email string, link string) (user entity.User, err error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	link = strings.TrimSpace(link)
	if name == "" || email == "" {
		return entity.User{}, fmt.Errorf("name and email are required")
	}
	if !utils.ValidateEmail(email) {
		return entity.User{}, fmt.Errorf("email is invalid")
	}
	if link != "" && !utils.ValidateURL(link) {
		link = ""
	}
	user = dao.FindUser(name, email)
	if user.IsEmpty() {
		user, err = dao.NewUser(name, email, link) // save a new user
		if err != nil {
			return entity.User{}, err
		}
	}
	return user, nil
}

func (dao *Dao) FindCreateNotify(userID uint, lookCommentID uint) entity.Notify {
	notify := dao.FindNotify(userID, lookCommentID)
	if notify.IsEmpty() {
		notify = dao.NewNotify(userID, lookCommentID)
	}
	return notify
}
