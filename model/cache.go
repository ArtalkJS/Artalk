package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

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

	entry, err := lib.CACHE.Get(lib.Ctx, name)
	if err != nil {
		return cacher, err
	}

	str := entry.([]byte)
	err = json.Unmarshal(str, destStruct)
	if err != nil {
		logrus.Debug("[缓存反序列化错误] ", name, " ", err)

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

	str, err := json.Marshal(srcStruct)
	if err != nil {
		return err
	}

	err = lib.CACHE.Set(lib.Ctx, name, []byte(str), &store.Options{})
	if err != nil {
		return err
	}

	logrus.Debug("[写入缓存] " + name)

	return nil
}

func ClearCache(name string) error {
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

func UserCacheClear(user *User) error {
	if err := ClearCache(fmt.Sprintf("user#id=%d", user.ID)); err != nil {
		return err
	}

	if err := ClearCache(fmt.Sprintf("user#name=%s;email=%s", strings.ToLower(user.Name), strings.ToLower(user.Email))); err != nil {
		return err
	}

	return nil
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

func SiteCacheClear(site *Site) error {
	if err := ClearCache(fmt.Sprintf("site#id=%d", site.ID)); err != nil {
		return err
	}

	if err := ClearCache(fmt.Sprintf("site#name=%s", site.Name)); err != nil {
		return err
	}

	return nil
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

func PageCacheClear(page *Page) error {
	if err := ClearCache(fmt.Sprintf("page#id=%d", page.ID)); err != nil {
		return err
	}

	if err := ClearCache(fmt.Sprintf("page#key=%s;site_name=%s", page.Key, page.SiteName)); err != nil {
		return err
	}

	return nil
}

func CommentCacheSave(comment *Comment) error {
	// 缓存 ID
	err := StoreCache(fmt.Sprintf("comment#id=%d", comment.ID), comment)
	if err != nil {
		return err
	}

	// 缓存 Rid
	if comment.Rid != 0 {
		_ChildCommentCacheSave(comment.Rid, comment.ID)
	}

	return err
}

func CommentCacheClear(comment *Comment) error {
	if err := ClearCache(fmt.Sprintf("comment#id=%d", comment.ID)); err != nil {
		return err
	}

	// 清除 Rid 缓存
	if comment.Rid != 0 {
		_ChildCommentCacheClear(comment.Rid, comment.ID)
	}

	return nil
}

// 缓存 父ID=>子ID 评论数据
func _ChildCommentCacheSave(parentID uint, childID uint) {
	var childIDs []uint
	cacher, err := FindCache(fmt.Sprintf("parent-comments#pid=%d", parentID), &childIDs)
	if err != nil {
		// 初始化
		childIDs = []uint{}
	}
	childIDs = append(childIDs, childID)
	cacher.StoreCache(func() interface{} {
		return &childIDs
	})
}

func _ChildCommentCacheClear(parentID uint, childID uint) {
	var childIDs []uint
	cacher, err := FindCache(fmt.Sprintf("parent-comments#pid=%d", parentID), &childIDs)
	if err != nil {
		return
	}

	// remove
	var nArr []uint
	for _, id := range childIDs {
		if id != childID {
			nArr = append(nArr, id)
		}
	}

	cacher.StoreCache(func() interface{} {
		return &nArr
	})
}
