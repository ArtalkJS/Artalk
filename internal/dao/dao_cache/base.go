package dao_cache

import (
	"fmt"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/cache"
	"github.com/ArtalkJS/Artalk/internal/entity"
)

type Cache struct {
}

func New() *Cache {
	return &Cache{}
}

// TODO remove cache functions away global states

func (c *Cache) UserCacheSave(user *entity.User) error {
	// 缓存 ID
	err := cache.StoreCache(fmt.Sprintf("user#id=%d", user.ID), user)
	if err != nil {
		return err
	}

	// 缓存 Name x Email
	err = cache.StoreCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)), user)
	if err != nil {
		return err
	}

	return err
}

func (c *Cache) UserCacheDel(user *entity.User) {
	cache.DelCache(fmt.Sprintf("user#id=%d", user.ID))
	cache.DelCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)))
}

func (c *Cache) SiteCacheSave(site *entity.Site) error {
	// 缓存 ID
	err := cache.StoreCache(fmt.Sprintf("site#id=%d", site.ID), site)
	if err != nil {
		return err
	}

	// 缓存 Name
	err = cache.StoreCache(fmt.Sprintf("site#name=%s", site.Name), site)
	if err != nil {
		return err
	}

	return err
}

func (c *Cache) SiteCacheDel(site *entity.Site) {
	cache.DelCache(fmt.Sprintf("site#id=%d", site.ID))
	cache.DelCache(fmt.Sprintf("site#name=%s", site.Name))
}

func (c *Cache) PageCacheSave(page *entity.Page) error {
	// 缓存 ID
	err := cache.StoreCache(fmt.Sprintf("page#id=%d", page.ID), page)
	if err != nil {
		return err
	}

	// 缓存 Key x SiteName
	err = cache.StoreCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName), page)
	if err != nil {
		return err
	}

	return err
}

func (c *Cache) PageCacheDel(page *entity.Page) {
	cache.DelCache(fmt.Sprintf("page#id=%d", page.ID))
	cache.DelCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName))
}

func (c *Cache) CommentCacheSave(comment *entity.Comment) error {
	// 缓存 ID
	err := cache.StoreCache(fmt.Sprintf("comment#id=%d", comment.ID), comment)
	if err != nil {
		return err
	}

	// 缓存 Rid
	if comment.Rid != 0 {
		c.ChildCommentCacheSave(comment.Rid, comment.ID)
	}

	return err
}

func (c *Cache) CommentCacheDel(comment *entity.Comment) {
	cache.DelCache(fmt.Sprintf("comment#id=%d", comment.ID))

	// 清除 Rid 缓存
	c.ChildCommentCacheDel(comment.ID)
	if comment.Rid != 0 {
		c.ChildCommentCacheDel(comment.Rid)
	}
}

// 缓存 父ID=>子ID 评论数据
func (c *Cache) ChildCommentCacheSave(parentID uint, childID uint) {
	var childIDs []uint
	var cacheName = fmt.Sprintf("parent-comments#pid=%d", parentID)

	err := cache.FindCache(cacheName, &childIDs)
	if err != nil { // 初始化
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

	cache.StoreCache(cacheName, childIDs)
}

func (c *Cache) ChildCommentCacheDel(parentID uint) {
	cache.DelCache(fmt.Sprintf("parent-comments#pid=%d", parentID))
}
