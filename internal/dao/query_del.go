package dao

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
)

// TODO: consider refactor make all delete operations in a transaction

func (dao *Dao) DelComment(comment *entity.Comment) error {
	// 清除 notify
	if err := dao.DB().Unscoped().Where("comment_id = ?", comment.ID).Delete(&entity.Notify{}).Error; err != nil {
		return err
	}

	// 清除 vote
	if err := dao.DB().Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		comment.ID,
		string(entity.VoteTypeCommentUp),
		string(entity.VoteTypeCommentDown),
	).Delete(&entity.Vote{}).Error; err != nil {
		return err
	}

	// 删除 comment
	err := dao.DB().Unscoped().Delete(comment).Error
	if err != nil {
		return err
	}

	// 删除缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.CommentCacheDel(comment)
	})

	return nil
}

// 删除所有子评论
func (dao *Dao) DelCommentChildren(parentID uint) error {
	var rErr error
	children := dao.FindCommentChildren(parentID)
	for _, c := range children {
		err := dao.DelComment(&c)
		if err != nil {
			rErr = err
		}
	}
	return rErr
}

func (dao *Dao) DelPage(page *entity.Page) error {
	err := dao.DB().Unscoped().Delete(page).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var comments []entity.Comment
	dao.DB().Where("page_key = ? AND site_name = ?", page.Key, page.SiteName).Find(&comments)

	for _, c := range comments {
		dao.DelComment(&c)
	}

	// 删除 vote
	dao.DB().Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		page.ID,
		string(entity.VoteTypePageUp),
		string(entity.VoteTypePageDown),
	).Delete(&entity.Vote{})

	// 删除缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.PageCacheDel(page)
	})

	return nil
}

func (dao *Dao) DelSite(site *entity.Site) error {
	err := dao.DB().Unscoped().Delete(&site).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var pages []entity.Page
	dao.DB().Where("site_name = ?", site.Name).Find(&pages)
	for _, p := range pages {
		dao.DelPage(&p)
	}

	// 删除缓存
	dao.CacheAction(func(cache *DaoCache) {
		cache.SiteCacheDel(site)
	})

	return nil
}

func (dao *Dao) DelUser(user *entity.User) error {
	err := dao.DB().Unscoped().Delete(&user).Error
	if err != nil {
		return err
	}

	// Delete user comments
	var comments []entity.Comment
	dao.DB().Where("user_id = ?", user.ID).Find(&comments)
	for _, c := range comments {
		dao.DelComment(&c)           // Delete parent comment
		dao.DelCommentChildren(c.ID) // Delete all child comments
	}

	// Delete user auth identities
	var authIdentity []entity.AuthIdentity
	dao.DB().Where("user_id = ?", user.ID).Find(&authIdentity)
	for _, a := range authIdentity {
		dao.DelAuthIdentity(&a)
	}

	// Clear cache
	dao.CacheAction(func(cache *DaoCache) {
		cache.UserCacheDel(user)
	})

	return nil
}

func (dao *Dao) DelAuthIdentity(authIdentity *entity.AuthIdentity) error {
	err := dao.DB().Unscoped().Delete(&authIdentity).Error
	if err != nil {
		return err
	}

	// TODO 删除缓存
	// dao.CacheAction(func(cache *DaoCache) {
	// 	cache.AuthIdentityCacheDel(authIdentity)
	// })

	return nil
}
