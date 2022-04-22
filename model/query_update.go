package model

import (
	"errors"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/sirupsen/logrus"
)

// 更新评论
func UpdateComment(comment *Comment) error {
	err := lib.DB.Save(comment).Error
	if err != nil {
		logrus.Error("Update Comment error: ", err)
	}
	// 更新缓存
	CommentCacheSave(comment)
	return err
}

func UpdateSite(site *Site) error {
	err := lib.DB.Save(site).Error
	if err != nil {
		logrus.Error("Update Site error: ", err)
	}
	SiteCacheSave(site)
	return err
}

func UpdateUser(user *User) error {
	err := lib.DB.Save(user).Error
	if err != nil {
		logrus.Error("Update User error: ", err)
	}
	UserCacheSave(user)
	return err
}

func UpdatePage(page *Page) error {
	err := lib.DB.Save(page).Error
	if err != nil {
		logrus.Error("Update Page error: ", err)
	}
	PageCacheSave(page)
	return err
}

func UserNotifyMarkAllAsRead(userID uint) error {
	if userID == 0 {
		return errors.New("user not found")
	}

	lib.DB.Model(&Notify{}).Where("user_id = ?", userID).Updates(&Notify{
		IsRead: true,
		ReadAt: time.Now(),
	})

	return nil
}
