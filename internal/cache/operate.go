package cache

import (
	"fmt"
	"time"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"gorm.io/gorm"
)

// 缓存预热
func CacheWarmUp(db *gorm.DB) {
	// Users
	{
		start := time.Now()

		var items []entity.User
		db.Find(&items)

		for _, item := range items {
			UserCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Users] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Sites
	{
		start := time.Now()

		var items []entity.Site
		db.Find(&items)

		for _, item := range items {
			SiteCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Sites] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Pages
	{
		start := time.Now()

		var items []entity.Page
		db.Find(&items)

		for _, item := range items {
			PageCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Pages] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Comments
	{
		start := time.Now()

		var items []entity.Comment
		db.Find(&items)

		for _, item := range items {
			CommentCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Comments] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}
}

// 清空缓存
func CacheFlushAll(db *gorm.DB) {
	// Users
	{
		var items []entity.User
		db.Find(&items)

		for _, item := range items {
			UserCacheDel(&item)
		}
	}

	// Sites
	{
		var items []entity.Site
		db.Find(&items)

		for _, item := range items {
			SiteCacheDel(&item)
		}
	}

	// Pages
	{
		var items []entity.Page
		db.Find(&items)

		for _, item := range items {
			PageCacheDel(&item)
		}
	}

	// Comments
	{
		var items []entity.Comment
		db.Find(&items)

		for _, item := range items {
			CommentCacheDel(&item)
		}
	}
}
