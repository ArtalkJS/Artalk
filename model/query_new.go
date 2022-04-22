package model

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/sirupsen/logrus"
)

func NewSite(name string, urls string) Site {
	site := Site{
		Name: name,
		Urls: urls,
	}

	err := CreateSite(&site)
	if err != nil {
		logrus.Error("Create Site error: ", err)
	}

	return site
}

func CreateSite(site *Site) error {
	err := lib.DB.Create(&site).Error
	if err != nil {
		return err
	}

	// 制备缓存
	SiteCacheSave(site)

	return nil
}

func NewUser(name string, email string, link string) User {
	user := User{
		Name:  name,
		Email: email,
		Link:  link,
	}

	err := CreateUser(&user)
	if err != nil {
		logrus.Error("Create User error: ", err)
	}

	return user
}

func CreateUser(user *User) error {
	err := lib.DB.Create(&user).Error
	if err != nil {
		return err
	}

	// 制备缓存
	UserCacheSave(user)

	return nil
}

func NewPage(key string, pageTitle string, siteName string) Page {
	page := Page{
		Key:      key,
		Title:    pageTitle,
		SiteName: siteName,
	}

	err := CreatePage(&page)
	if err != nil {
		logrus.Error("Create Page error: ", err)
	}

	return page
}

func CreatePage(page *Page) error {
	err := lib.DB.Create(&page).Error
	if err != nil {
		return err
	}

	// 制备缓存
	PageCacheSave(page)

	return nil
}

func CreateComment(comment *Comment) error {
	err := lib.DB.Create(&comment).Error
	if err != nil {
		return err
	}

	// 制备缓存
	CommentCacheSave(comment)

	return nil
}

func NewNotify(userID uint, commentID uint) Notify {
	notify := Notify{
		UserID:    userID,
		CommentID: commentID,
		IsRead:    false,
		IsEmailed: false,
	}
	notify.GenerateKey()

	err := lib.DB.Create(&notify).Error
	if err != nil {
		logrus.Error("Create Notify error: ", err)
	}

	return notify
}
