package model

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func FindComment(id uint, siteID uint) Comment {
	var comment Comment
	lib.DB.Preload(clause.Associations).Where("id = ? AND site_id = ?", id, siteID).First(&comment)
	return comment
}

// 查找用户（返回：精确查找 AND）
func FindUser(name string, email string) User {
	var user User
	lib.DB.Where(&User{Name: name, Email: email}).First(&user)
	return user
}

func IsAdminUser(name string, email string) bool {
	var user User
	lib.DB.Where(&User{Name: name, Email: email, IsAdmin: true}).First(&user)
	return !user.IsEmpty()
}

func UpdateComment(comment *Comment) error {
	err := lib.DB.Save(comment).Error
	if err != nil {
		logrus.Error("Update Comment error: ", err)
	}
	return err
}

func FindSite(name string) Site {
	var site Site
	lib.DB.Where(&Site{Name: name}).First(&site)
	return site
}

func FindCreatePage(pageKey string, pageUrl string, pageTitle string, siteID uint) Page {
	page := FindPage(pageKey, siteID)
	if page.IsEmpty() {
		page = NewPage(pageKey, pageUrl, pageTitle, siteID)
	}
	return page
}

func FindCreateUser(name string, email string) User {
	user := FindUser(name, email)
	if user.IsEmpty() {
		user = NewUser(name, email) // save a new user
	}
	return user
}

func NewUser(name string, email string) User {
	user := User{
		Name:  name,
		Email: email,
	}

	err := lib.DB.Create(&user).Error
	if err != nil {
		logrus.Error("Save User error: ", err)
	}

	return user
}

func UpdateUser(user *User) error {
	err := lib.DB.Save(user).Error
	if err != nil {
		logrus.Error("Update User error: ", err)
	}

	return err
}

func FindPage(key string, siteID uint) Page {
	var page Page
	lib.DB.Where(&Page{Key: key, SiteID: siteID}).First(&page)
	return page
}

func FindPageByID(id uint, siteID uint) Page {
	var page Page
	lib.DB.Where("id = ? AND site_id = ?", id, siteID).First(&page)
	return page
}

func NewPage(key string, pageUrl string, pageTitle string, siteID uint) Page {
	page := Page{
		Key:    key,
		Url:    pageUrl,
		Title:  pageTitle,
		SiteID: siteID,
	}

	err := lib.DB.Create(&page).Error
	if err != nil {
		logrus.Error("Save Page error: ", err)
	}

	return page
}

func UpdatePage(page *Page) error {
	err := lib.DB.Save(page).Error
	if err != nil {
		logrus.Error("Update Page error: ", err)
	}
	return err
}
