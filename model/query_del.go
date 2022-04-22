package model

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
)

func DelComment(commentID uint) error {
	// 查询 Comment
	comment := FindComment(commentID)

	// 清除 notify
	if err := lib.DB.Unscoped().Where("comment_id = ?", commentID).Delete(&Notify{}).Error; err != nil {
		return err
	}

	// 清除 vote
	if err := lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		commentID,
		string(VoteTypeCommentUp),
		string(VoteTypeCommentDown),
	).Delete(&Vote{}).Error; err != nil {
		return err
	}

	// 删除 comment
	err := lib.DB.Unscoped().Delete(&Comment{}, commentID).Error
	if err != nil {
		return err
	}

	// 删除缓存
	CommentCacheClear(&comment)

	return nil
}

func DelPage(page *Page) error {
	err := lib.DB.Unscoped().Delete(page).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var comments []Comment
	lib.DB.Where("page_key = ? AND site_name = ?", page.Key, page.SiteName).Find(&comments)

	for _, c := range comments {
		DelComment(c.ID)
	}

	// 删除 vote
	lib.DB.Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		page.ID,
		string(VoteTypePageUp),
		string(VoteTypePageDown),
	).Delete(&Vote{})

	// 删除缓存
	PageCacheClear(page)

	return nil
}

func DelSite(site *Site, softDel bool) error {
	err := lib.DB.Unscoped().Delete(&site).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	if !softDel {
		var pages []Page
		lib.DB.Where("site_name = ?", site.Name).Find(&pages)

		for _, p := range pages {
			DelPage(&p)
		}
	}

	// 删除缓存
	SiteCacheClear(site)

	return nil
}

func DelUser(user *User) error {
	err := lib.DB.Unscoped().Delete(&user).Error
	if err != nil {
		return err
	}

	// 删除缓存
	UserCacheClear(user)

	return nil
}
