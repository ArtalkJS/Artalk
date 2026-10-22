package dao

import (
	"fmt"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"gorm.io/gorm"
)

// 更新评论
func (dao *Dao) UpdateComment(comment *entity.Comment) error {
	err := dao.DB().Save(comment).Error
	if err != nil {
		log.Error("Update Comment error: ", err)
	}
	// 更新缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.CommentCacheSave(comment)
	})
	return err
}

func (dao *Dao) UpdateSite(site *entity.Site) error {
	err := dao.DB().Save(site).Error
	if err != nil {
		log.Error("Update Site error: ", err)
	}
	dao.CacheAction(func(cache *DaoCache) {
		cache.SiteCacheSave(site)
	})
	return err
}

func (dao *Dao) UpdateUser(user *entity.User) error {
	err := dao.DB().Save(user).Error
	if err != nil {
		log.Error("Update User error: ", err)
	}
	dao.CacheAction(func(cache *DaoCache) {
		cache.UserCacheSave(user)
	})
	return err
}

func (dao *Dao) UpdatePage(page *entity.Page) error {
	err := dao.DB().Save(page).Error
	if err != nil {
		log.Error("Update Page error: ", err)
	}
	dao.CacheAction(func(cache *DaoCache) {
		cache.PageCacheSave(page)
	})
	return err
}

// IncrementPagePV does not write a potentially stale page snapshot back to the
// database. Atomic increments also preserve views from other server instances.
func (dao *Dao) IncrementPagePV(page *entity.Page) error {
	result := dao.DB().Model(&entity.Page{}).Where("id = ?", page.ID).
		Update("pv", gorm.Expr("COALESCE(pv, 0) + ?", 1))
	dao.CacheAction(func(cache *DaoCache) { cache.PageCacheDel(page) })
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return dao.DB().First(page, page.ID).Error
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

func (dao *Dao) UpdateAuthIdentity(authIdentity *entity.AuthIdentity) error {
	err := dao.DB().Save(authIdentity).Error
	if err != nil {
		log.Error("Update AuthIdentity error: ", err)
	}
	// TODO: 更新缓存
	// dao.CacheAction(func(cache *DaoCache) {
	// 	cache.AuthIdentityCacheSave(authIdentity)
	// })
	return err
}
