package model

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/eko/gocache/v2/store"
	"github.com/sirupsen/logrus"
)

var (
	MutexCache = sync.Mutex{}
)

type cacher struct{ cacheKey string }

func (c *cacher) StoreCache(getSrcStruct func() interface{}) error {
	return StoreCache(c.cacheKey, nil, getSrcStruct)
}

func FindCache(name string, destStruct interface{}) (cacher, error) {
	cacher := cacher{cacheKey: name}

	if !config.Instance.Cache.Enabled {
		return cacher, errors.New("缓存功能禁用")
	}

	_, err := lib.CACHE_marshal.Get(lib.Ctx, name, destStruct)
	if err != nil {
		return cacher, err
	}

	logrus.Debug("[缓存命中] ", name)

	return cacher, nil
}

func StoreCache(name string, srcStruct interface{}, getSrcStruct ...func() interface{}) error {
	MutexCache.Lock()
	defer MutexCache.Unlock()

	if len(getSrcStruct) > 0 { // getSrcStruct 为可选参数，当存在时会覆盖 srcStruct 参数
		srcStruct = getSrcStruct[0]() // 这个 func 内再执行 db 查询，加锁防止反复查询
	}

	if !config.Instance.Cache.Enabled {
		return nil
	}

	err := lib.CACHE_marshal.Set(lib.Ctx, name, srcStruct, &store.Options{
		Expiration: time.Duration(config.Instance.Cache.GetExpiresTime()),
	})
	if err != nil {
		return err
	}

	logrus.Debug("[写入缓存] " + name)

	return nil
}

func DelCache(name string) error {
	if !config.Instance.Cache.Enabled {
		return nil
	}

	MutexCache.Lock()
	defer MutexCache.Unlock()

	return lib.CACHE.Delete(lib.Ctx, name)
}

func UserCacheSave(user *User) error {
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

func UserCacheDel(user *User) {
	DelCache(fmt.Sprintf("user#id=%d", user.ID))
	DelCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email)))
}

func SiteCacheSave(site *Site) error {
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

func SiteCacheDel(site *Site) {
	DelCache(fmt.Sprintf("site#id=%d", site.ID))
	DelCache(fmt.Sprintf("site#name=%s", site.Name))
}

func PageCacheSave(page *Page) error {
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

func PageCacheDel(page *Page) {
	DelCache(fmt.Sprintf("page#id=%d", page.ID))
	DelCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName))
}

func CommentCacheSave(comment *Comment) error {
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

func CommentCacheDel(comment *Comment) {
	DelCache(fmt.Sprintf("comment#id=%d", comment.ID))

	// 清除 Rid 缓存
	ChildCommentCacheDel(comment.ID)
}

// 缓存 父ID=>子ID 评论数据
func ChildCommentCacheSave(parentID uint, childID uint) {
	var cacheKey = fmt.Sprintf("parent-comments#pid=%d", parentID)
	var childIDs []uint
	StoreCache(cacheKey, nil, func() interface{} {
		_, err := FindCache(cacheKey, &childIDs)
		if err != nil {
			// 初始化
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

		return &childIDs
	})
}

func ChildCommentCacheDel(parentID uint) {
	DelCache(fmt.Sprintf("parent-comments#pid=%d", parentID))
}

// 仅从 pid => child_ids 切片中删除一项，
// 最后结果可能出现：这个切片为空，但缓存存在 (可命中) 的情况
func ChildCommentCacheSplice(parentID uint, childID uint) {
	cacheKey := fmt.Sprintf("parent-comments#pid=%d", parentID)
	var childIDs []uint
	StoreCache(cacheKey, nil, func() interface{} {
		_, err := FindCache(cacheKey, &childIDs)
		if err != nil {
			return []uint{}
		}

		// remove item
		var nArr []uint
		for _, id := range childIDs {
			if id != childID {
				nArr = append(nArr, id)
			}
		}

		return &nArr
	})
}
