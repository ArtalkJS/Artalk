package query

import (
	"errors"
	"time"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/sirupsen/logrus"
)

// 更新评论
func UpdateComment(comment *entity.Comment) error {
	err := DB().Save(comment).Error
	if err != nil {
		logrus.Error("Update Comment error: ", err)
	}
	// 更新缓存
	cache.CommentCacheSave(comment)
	return err
}

func UpdateSite(site *entity.Site) error {
	err := DB().Save(site).Error
	if err != nil {
		logrus.Error("Update Site error: ", err)
	}
	cache.SiteCacheSave(site)
	return err
}

func UpdateUser(user *entity.User) error {
	err := DB().Save(user).Error
	if err != nil {
		logrus.Error("Update User error: ", err)
	}
	cache.UserCacheSave(user)
	return err
}

func UpdatePage(page *entity.Page) error {
	err := DB().Save(page).Error
	if err != nil {
		logrus.Error("Update Page error: ", err)
	}
	cache.PageCacheSave(page)
	return err
}

func UserNotifyMarkAllAsRead(userID uint) error {
	if userID == 0 {
		return errors.New("user not found")
	}

	nowTime := time.Now()

	DB().Model(&entity.Notify{}).Where("user_id = ?", userID).Updates(&entity.Notify{
		IsRead: true,
		ReadAt: &nowTime,
	})

	return nil
}
