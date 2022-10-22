package model

func DelComment(comment *Comment) error {
	// 清除 notify
	if err := DB().Unscoped().Where("comment_id = ?", comment.ID).Delete(&Notify{}).Error; err != nil {
		return err
	}

	// 清除 vote
	if err := DB().Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		comment.ID,
		string(VoteTypeCommentUp),
		string(VoteTypeCommentDown),
	).Delete(&Vote{}).Error; err != nil {
		return err
	}

	// 删除 comment
	err := DB().Unscoped().Delete(comment).Error
	if err != nil {
		return err
	}

	// 删除缓存
	CommentCacheDel(comment)

	return nil
}

// 删除所有子评论
func DelCommentChildren(parentID uint) error {
	var rErr error
	children := FindCommentChildren(parentID)
	for _, c := range children {
		err := DelComment(&c)
		if err != nil {
			rErr = err
		}
	}
	return rErr
}

func DelPage(page *Page) error {
	err := DB().Unscoped().Delete(page).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var comments []Comment
	DB().Where("page_key = ? AND site_name = ?", page.Key, page.SiteName).Find(&comments)

	for _, c := range comments {
		DelComment(&c)
	}

	// 删除 vote
	DB().Unscoped().Where(
		"target_id = ? AND (type = ? OR type = ?)",
		page.ID,
		string(VoteTypePageUp),
		string(VoteTypePageDown),
	).Delete(&Vote{})

	// 删除缓存
	PageCacheDel(page)

	return nil
}

func DelSite(site *Site) error {
	err := DB().Unscoped().Delete(&site).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var pages []Page
	DB().Where("site_name = ?", site.Name).Find(&pages)
	for _, p := range pages {
		DelPage(&p)
	}

	// 删除缓存
	SiteCacheDel(site)

	return nil
}

func DelUser(user *User) error {
	err := DB().Unscoped().Delete(&user).Error
	if err != nil {
		return err
	}

	// 删除所有相关内容
	var comments []Comment
	DB().Where("user_id = ?", user.ID).Find(&comments)
	for _, c := range comments {
		DelComment(&c)           // 删除主评论
		DelCommentChildren(c.ID) // 删除子评论
	}

	// 删除缓存
	UserCacheDel(user)

	return nil
}
