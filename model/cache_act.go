package model

import (
	"fmt"
	"time"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/sirupsen/logrus"
)

// 缓存预热
func CacheWarmUp() {
	// Users
	{
		start := time.Now()

		var items []User
		lib.DB.Find(&items)

		for _, item := range items {
			UserCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Users] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Sites
	{
		start := time.Now()

		var items []Site
		lib.DB.Find(&items)

		for _, item := range items {
			SiteCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Sites] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Pages
	{
		start := time.Now()

		var items []Page
		lib.DB.Find(&items)

		for _, item := range items {
			PageCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Pages] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}

	// Comments
	{
		start := time.Now()

		var items []Comment
		lib.DB.Find(&items)

		for _, item := range items {
			CommentCacheSave(&item)
		}

		logrus.Debug(fmt.Sprintf("[Comments] 缓存完毕 (共 %d 个，耗时：%s)", len(items), time.Since(start)))
	}
}

// 清空缓存
func CacheFlushAll() {
	// Users
	{
		var items []User
		lib.DB.Find(&items)

		for _, item := range items {
			UserCacheDel(&item)
		}
	}

	// Sites
	{
		var items []Site
		lib.DB.Find(&items)

		for _, item := range items {
			SiteCacheDel(&item)
		}
	}

	// Pages
	{
		var items []Page
		lib.DB.Find(&items)

		for _, item := range items {
			PageCacheDel(&item)
		}
	}

	// Comments
	{
		var items []Comment
		lib.DB.Find(&items)

		for _, item := range items {
			CommentCacheDel(&item)
		}
	}
}
