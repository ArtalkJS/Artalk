package dao

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/entity"
)

type DaoCache struct {
	*cache.Cache
}

func NewCacheAdaptor(cache *cache.Cache) *DaoCache {
	return &DaoCache{Cache: cache}
}

func (c *DaoCache) UserCacheSave(user *entity.User) error {
	// 缓存 ID
	err := c.StoreCache(fmt.Sprintf("user#id=%d", user.ID), user)
	if err != nil {
		return err
	}

	// 缓存 Name x Email
	err = c.StoreCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)), user)
	if err != nil {
		return err
	}

	return err
}

func (c *DaoCache) UserCacheDel(user *entity.User) {
	c.DelCache(fmt.Sprintf("user#id=%d", user.ID))
	c.DelCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)))
}

func (c *DaoCache) SiteCacheSave(site *entity.Site) error {
	// 缓存 ID
	err := c.StoreCache(fmt.Sprintf("site#id=%d", site.ID), site)
	if err != nil {
		return err
	}

	// 缓存 Name
	err = c.StoreCache(fmt.Sprintf("site#name=%s", site.Name), site)
	if err != nil {
		return err
	}

	return err
}

func (c *DaoCache) SiteCacheDel(site *entity.Site) {
	c.DelCache(fmt.Sprintf("site#id=%d", site.ID))
	c.DelCache(fmt.Sprintf("site#name=%s", site.Name))
}

func (c *DaoCache) PageCacheSave(page *entity.Page) error {
	// 缓存 ID
	err := c.StoreCache(fmt.Sprintf("page#id=%d", page.ID), page)
	if err != nil {
		return err
	}

	// 缓存 Key x SiteName
	err = c.StoreCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName), page)
	if err != nil {
		return err
	}

	return err
}

func (c *DaoCache) PageCacheDel(page *entity.Page) {
	c.DelCache(fmt.Sprintf("page#id=%d", page.ID))
	c.DelCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName))
}

func (c *DaoCache) CommentCacheSave(comment *entity.Comment) error {
	// 缓存 ID
	err := c.StoreCache(fmt.Sprintf("comment#id=%d", comment.ID), comment)
	if err != nil {
		return err
	}

	// 缓存 Rid
	if comment.Rid != 0 {
		c.ChildCommentCacheSave(comment.Rid, comment.ID)
	}

	return err
}

func (c *DaoCache) CommentCacheDel(comment *entity.Comment) {
	c.DelCache(fmt.Sprintf("comment#id=%d", comment.ID))

	// 清除 Rid 缓存
	c.ChildCommentCacheDel(comment.ID)
	if comment.Rid != 0 {
		c.ChildCommentCacheDel(comment.Rid)
	}
}

// 缓存 父ID=>子ID 评论数据
func (c *DaoCache) ChildCommentCacheSave(parentID uint, childID uint) {
	var cacheName = fmt.Sprintf("parent-comments#pid=%d", parentID)

	var childIDs []uint
	err := c.FindCache(cacheName, &childIDs)
	if err != nil || childIDs == nil {
		childIDs = []uint{}
	}

	isExist := false
	for _, i := range childIDs {
		if i == childID {
			isExist = true
			break
		}
	}

	// append if Not exist
	if !isExist {
		childIDs = append(childIDs, childID)
	}

	c.StoreCache(cacheName, childIDs)
}

func (c *DaoCache) ChildCommentCacheDel(parentID uint) {
	c.DelCache(fmt.Sprintf("parent-comments#pid=%d", parentID))
}
