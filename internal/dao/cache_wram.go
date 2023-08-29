package dao

import (
	"fmt"
	"time"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

// 缓存预热
func (dao *Dao) CacheWarmUp() {
	// Users
	{
		start := time.Now()

		var items []entity.User
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.UserCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Users] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Sites
	{
		start := time.Now()

		var items []entity.Site
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.SiteCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Sites] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Pages
	{
		start := time.Now()

		var items []entity.Page
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.PageCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Pages] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Comments
	{
		start := time.Now()

		var items []entity.Comment
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.CommentCacheSave(&item)
		}

		log.Debug(fmt.Sprintf("[Comments] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}
}

// 清空缓存
func (dao *Dao) CacheFlushAll() {
	// Users
	{
		var items []entity.User
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.UserCacheDel(&item)
		}
	}

	// Sites
	{
		var items []entity.Site
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.SiteCacheDel(&item)
		}
	}

	// Pages
	{
		var items []entity.Page
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.PageCacheDel(&item)
		}
	}

	// Comments
	{
		var items []entity.Comment
		dao.DB().Find(&items)

		for _, item := range items {
			dao.cache.CommentCacheDel(&item)
		}
	}
}
