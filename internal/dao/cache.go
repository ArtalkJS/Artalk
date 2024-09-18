package dao

import (
	"fmt"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/cache"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"golang.org/x/exp/slices"
)

const (
	UserByIDKey            = "user#id=%d"
	UserByNameEmailKey     = "user#name=%s;email=%s"
	UserIDByEmailKey       = "user_id#email=%s"
	SiteByIDKey            = "site#id=%d"
	SiteByNameKey          = "site#name=%s"
	PageByIDKey            = "page#id=%d"
	PageByKeySiteNameKey   = "page#key=%s;site_name=%s"
	CommentByIDKey         = "comment#id=%d"
	CommentChildIDsByIDKey = "comment_child_ids#id=%d"
	NotifyByUserCommentKey = "notify#user_id=%d;comment_id=%d"
)

type DaoCache struct {
	*cache.Cache
}

func NewCacheAdaptor(cache *cache.Cache) *DaoCache {
	return &DaoCache{Cache: cache}
}

func (c *DaoCache) UserCacheSave(user *entity.User) error {
	return c.StoreCache(user,
		fmt.Sprintf(UserByIDKey, user.ID),
		fmt.Sprintf(UserByNameEmailKey, strings.ToLower(user.Name), strings.ToLower(user.Email)),
	)
}

func (c *DaoCache) UserCacheDel(user *entity.User) {
	c.DelCache(
		fmt.Sprintf(UserByIDKey, user.ID),
		fmt.Sprintf(UserByNameEmailKey, strings.ToLower(user.Name), strings.ToLower(user.Email)),
		fmt.Sprintf(UserIDByEmailKey, strings.ToLower(user.Email)),
	)
}

func (c *DaoCache) SiteCacheSave(site *entity.Site) error {
	return c.StoreCache(site,
		fmt.Sprintf(SiteByIDKey, site.ID),
		fmt.Sprintf(SiteByNameKey, site.Name),
	)
}

func (c *DaoCache) SiteCacheDel(site *entity.Site) {
	c.DelCache(fmt.Sprintf(SiteByIDKey, site.ID))
	c.DelCache(fmt.Sprintf(SiteByNameKey, site.Name))
}

func (c *DaoCache) PageCacheSave(page *entity.Page) error {
	return c.StoreCache(page,
		fmt.Sprintf(PageByIDKey, page.ID),
		fmt.Sprintf(PageByKeySiteNameKey, page.Key, page.SiteName),
	)
}

func (c *DaoCache) PageCacheDel(page *entity.Page) {
	c.DelCache(
		fmt.Sprintf(PageByIDKey, page.ID),
		fmt.Sprintf(PageByKeySiteNameKey, page.Key, page.SiteName),
	)
}

func (c *DaoCache) CommentCacheSave(comment *entity.Comment) (err error) {
	if storeErr := c.StoreCache(comment, fmt.Sprintf(CommentByIDKey, comment.ID)); storeErr != nil {
		err = storeErr
	}
	if storeErr := c.ChildCommentCacheSave(comment); storeErr != nil {
		err = storeErr
	}
	return
}

func (c *DaoCache) CommentCacheDel(comment *entity.Comment) {
	c.DelCache(fmt.Sprintf(CommentByIDKey, comment.ID))
	c.ChildCommentCacheDel(comment)
}

// 缓存 父ID=>子ID 评论数据
func (c *DaoCache) ChildCommentCacheSave(comment *entity.Comment) error {
	// 若 comment 为根评论
	if comment.Rid == 0 {
		// 查询是否有缓存数据
		// 若无缓存数据，则创建初始空数据，无子评论
		// 若有缓存数据，说明 comment 是 update 操作，则不更改
		cacheName := fmt.Sprintf(CommentChildIDsByIDKey, comment.ID)
		var childIDs []uint
		if err := c.FindCache(cacheName, &childIDs); err != nil {
			// 无缓存，则初始化
			return c.StoreCache([]uint{}, cacheName)
		}

		return nil
	}

	// 若 comment 为子评论（Rid 不为空，Rid 为父 ID），
	// 则为父评论的 “子评论 ID 列表” 追加 “子评论 ID”
	parentID := comment.Rid
	childID := comment.ID

	var cacheName = fmt.Sprintf(CommentChildIDsByIDKey, parentID)

	var childIDs []uint
	c.FindCache(cacheName, &childIDs)

	// append if childID Not contains
	if !slices.Contains(childIDs, childID) {
		childIDs = append(childIDs, childID)
	}

	return c.StoreCache(childIDs, cacheName)
}

func (c *DaoCache) ChildCommentCacheDel(comment *entity.Comment) {
	c.DelCache(fmt.Sprintf(CommentChildIDsByIDKey, comment.ID))
	if comment.Rid != 0 {
		c.DelCache(fmt.Sprintf(CommentChildIDsByIDKey, comment.Rid))
		// TODO 从 slice 中 remove 并更新缓存而不是直接删除缓存，删除缓存会导致重新查询 db 而性能下降
	}
}
