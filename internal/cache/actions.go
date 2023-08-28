package cache

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/eko/gocache/lib/v4/store"
	"golang.org/x/sync/singleflight"
)

var (
	CacheFindGroup = new(singleflight.Group)
)

func FindAndStoreCache(name string, dest interface{}, queryDBResult func() interface{}) error {
	// SingleFlight 防止缓存击穿 (Cache breakdown)
	v, err, _ := CacheFindGroup.Do(name, func() (interface{}, error) {
		err := FindCache(name, dest)

		// cache hit 直接返回结果
		if err == nil {
			return dest, nil
		}

		// cache miss 查数据库
		result := queryDBResult()
		if err := StoreCache(name, result); err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return err
	}

	if v != nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v).Elem()) // similar to `*dest = &v`
	}

	return nil
}

func FindCache(name string, dest interface{}) error {
	if std == nil {
		return fmt.Errorf("cache not initialize")
	}

	// `Get()` is Thread Safe, so no need to add Mutex
	// @see https://github.com/go-redis/redis/issues/23
	_, err := std.Get(ctx, name, dest)
	if err != nil {
		return err
	}

	log.Debug("[Cache Hit] " + name)

	return nil
}

func StoreCache(name string, source interface{}) error {
	if std == nil {
		return fmt.Errorf("cache not initialize")
	}

	// `Set()` is Thread Safe too, no need to add Mutex either
	err := std.Set(ctx, name, source,
		store.WithExpiration(ttl),
	)
	if err != nil {
		return err
	}

	log.Debug("[写入缓存] " + name)

	return nil
}

func DelCache(name string) error {
	if std == nil {
		return fmt.Errorf("cache not initialize")
	}

	return std.Delete(ctx, name)
}

func UserCacheSave(user *entity.User) error {
	// 缓存 ID
	err := StoreCache(fmt.Sprintf("user#id=%d", user.ID), user)
	if err != nil {
		return err
	}

	// 缓存 Name x Email
	err = StoreCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)), user)
	if err != nil {
		return err
	}

	return err
}

func UserCacheDel(user *entity.User) {
	DelCache(fmt.Sprintf("user#id=%d", user.ID))
	DelCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)))
}

func SiteCacheSave(site *entity.Site) error {
	// 缓存 ID
	err := StoreCache(fmt.Sprintf("site#id=%d", site.ID), site)
	if err != nil {
		return err
	}

	// 缓存 Name
	err = StoreCache(fmt.Sprintf("site#name=%s", site.Name), site)
	if err != nil {
		return err
	}

	return err
}

func SiteCacheDel(site *entity.Site) {
	DelCache(fmt.Sprintf("site#id=%d", site.ID))
	DelCache(fmt.Sprintf("site#name=%s", site.Name))
}

func PageCacheSave(page *entity.Page) error {
	// 缓存 ID
	err := StoreCache(fmt.Sprintf("page#id=%d", page.ID), page)
	if err != nil {
		return err
	}

	// 缓存 Key x SiteName
	err = StoreCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName), page)
	if err != nil {
		return err
	}

	return err
}

func PageCacheDel(page *entity.Page) {
	DelCache(fmt.Sprintf("page#id=%d", page.ID))
	DelCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName))
}

func CommentCacheSave(comment *entity.Comment) error {
	// 缓存 ID
	err := StoreCache(fmt.Sprintf("comment#id=%d", comment.ID), comment)
	if err != nil {
		return err
	}

	// 缓存 Rid
	if comment.Rid != 0 {
		ChildCommentCacheSave(comment.Rid, comment.ID)
	}

	return err
}

func CommentCacheDel(comment *entity.Comment) {
	DelCache(fmt.Sprintf("comment#id=%d", comment.ID))

	// 清除 Rid 缓存
	ChildCommentCacheDel(comment.ID)
	if comment.Rid != 0 {
		ChildCommentCacheDel(comment.Rid)
	}
}

// 缓存 父ID=>子ID 评论数据
func ChildCommentCacheSave(parentID uint, childID uint) {
	var childIDs []uint
	var cacheName = fmt.Sprintf("parent-comments#pid=%d", parentID)

	err := FindCache(cacheName, &childIDs)
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

	StoreCache(cacheName, childIDs)
}

func ChildCommentCacheDel(parentID uint) {
	DelCache(fmt.Sprintf("parent-comments#pid=%d", parentID))
}
