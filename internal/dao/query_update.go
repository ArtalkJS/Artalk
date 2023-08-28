package dao

import (
	"fmt"
	"time"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

// 更新评论
func (dao *Dao) UpdateComment(comment *entity.Comment) error {
	err := dao.DB().Save(comment).Error
	if err != nil {
		log.Error("Update Comment error: ", err)
	}
	// 更新缓存
	cache.CommentCacheSave(comment)
	return err
}

func (dao *Dao) UpdateSite(site *entity.Site) error {
	err := dao.DB().Save(site).Error
	if err != nil {
		log.Error("Update Site error: ", err)
	}
	cache.SiteCacheSave(site)
	return err
}

func (dao *Dao) UpdateUser(user *entity.User) error {
	err := dao.DB().Save(user).Error
	if err != nil {
		log.Error("Update User error: ", err)
	}
	cache.UserCacheSave(user)
	return err
}

func (dao *Dao) UpdatePage(page *entity.Page) error {
	err := dao.DB().Save(page).Error
	if err != nil {
		log.Error("Update Page error: ", err)
	}
	cache.PageCacheSave(page)
	return err
}

func (dao *Dao) UserNotifyMarkAllAsRead(userID uint) error {
	if userID == 0 {
		return fmt.Errorf("user not found")
	}

	nowTime := time.Now()

	dao.DB().Model(&entity.Notify{}).Where("user_id = ?", userID).Updates(&entity.Notify{
		IsRead: true,
		ReadAt: &nowTime,
	})

	return nil
}
