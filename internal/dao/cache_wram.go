package dao

import (
	"fmt"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
)

// 缓存预热
func (dao *Dao) CacheWarmUp() {
	// Users
	{
		start := time.Now()

		var items []entity.User
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.UserCacheSave(&item)
			})
		}

		log.Debug(fmt.Sprintf("[Users] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Sites
	{
		start := time.Now()

		var items []entity.Site
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.SiteCacheSave(&item)
			})
		}

		log.Debug(fmt.Sprintf("[Sites] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Pages
	{
		start := time.Now()

		var items []entity.Page
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.PageCacheSave(&item)
			})
		}

		log.Debug(fmt.Sprintf("[Pages] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Comments
	{
		start := time.Now()

		var items []entity.Comment
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.CommentCacheSave(&item)

			})
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
			dao.CacheAction(func(cache *DaoCache) {
				cache.UserCacheDel(&item)
			})
		}
	}

	// Sites
	{
		var items []entity.Site
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.SiteCacheDel(&item)
			})
		}
	}

	// Pages
	{
		var items []entity.Page
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.PageCacheDel(&item)
			})
		}
	}

	// Comments
	{
		var items []entity.Comment
		dao.DB().Find(&items)

		for _, item := range items {
			dao.CacheAction(func(cache *DaoCache) {
				cache.CommentCacheDel(&item)
			})
		}
	}
}
