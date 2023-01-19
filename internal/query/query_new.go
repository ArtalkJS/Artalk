package query

import (
	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/sirupsen/logrus"
)

func NewSite(name string, urls string) entity.Site {
	site := entity.Site{
		Name: name,
		Urls: urls,
	}

	err := CreateSite(&site)
	if err != nil {
		logrus.Error("Create Site error: ", err)
	}

	return site
}

func CreateSite(site *entity.Site) error {
	err := DB().Create(&site).Error
	if err != nil {
		return err
	}

	// 制备缓存
	cache.SiteCacheSave(site)

	return nil
}

func NewUser(name string, email string, link string) entity.User {
	user := entity.User{
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

func CreateUser(user *entity.User) error {
	err := DB().Create(&user).Error
	if err != nil {
		return err
	}

	// 制备缓存
	cache.UserCacheSave(user)

	return nil
}

func NewPage(key string, pageTitle string, siteName string) entity.Page {
	page := entity.Page{
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

func CreatePage(page *entity.Page) error {
	err := DB().Create(&page).Error
	if err != nil {
		return err
	}

	// 制备缓存
	cache.PageCacheSave(page)

	return nil
}

func CreateComment(comment *entity.Comment) error {
	err := DB().Create(&comment).Error
	if err != nil {
		return err
	}

	// 制备缓存
	cache.CommentCacheSave(comment)

	return nil
}

func NewNotify(userID uint, commentID uint) entity.Notify {
	notify := entity.Notify{
		UserID:    userID,
		CommentID: commentID,
		IsRead:    false,
		IsEmailed: false,
	}
	notify.GenerateKey()

	err := DB().Create(&notify).Error
	if err != nil {
		logrus.Error("Create Notify error: ", err)
	}

	return notify
}

func NewVote(targetID uint, voteType entity.VoteType, userID uint, ua string, ip string) (entity.Vote, error) {
	vote := entity.Vote{
		TargetID: targetID,
		Type:     voteType,
		UserID:   userID,
		UA:       ua,
		IP:       ip,
	}

	err := DB().Create(&vote).Error
	if err != nil {
		logrus.Error("Create Vote error: ", err)
	}

	return vote, err
}
