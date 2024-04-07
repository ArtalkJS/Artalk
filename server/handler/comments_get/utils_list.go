package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// Find all nested children (for nested mode)
func findNestedChildren(dao *dao.Dao, comments []entity.CookedComment, commonScopes []func(*gorm.DB) *gorm.DB) []entity.CookedComment {
	allRootIDs := lo.Map(comments, func(c entity.CookedComment, _ int) uint { return c.ID })
	// TODO: Add pagination for nested mode
	// 	All children will be loaded at once without pagination, which may cause performance issues.
	// 	The backend will response all to the client-side, and render by the client-side itself.
	var children []*entity.Comment
	dao.DB().Model(&entity.Comment{}).
		Scopes(commonScopes...).
		Where("root_id IN ?", allRootIDs).
		Find(&children)
	comments = append(comments, dao.CookAllComments(children)...)
	return comments
}

// Find all linked comments (for flat mode)
func findFlatLinkedComments(dao *dao.Dao, comments []entity.CookedComment, commonScopes []func(*gorm.DB) *gorm.DB) []entity.CookedComment {
	allCommentIDs := map[uint]bool{}
	for _, c := range comments {
		allCommentIDs[c.ID] = true
	}

	missCommentIDs := lo.Map(lo.Filter(comments, func(c entity.CookedComment, i int) bool {
		return c.Rid != 0 && !allCommentIDs[c.Rid]
	}), func(c entity.CookedComment, i int) uint {
		return c.Rid
	})

	// Find linked comments which `id` is the same as `rid`
	var linkedComments []entity.Comment
	if len(missCommentIDs) > 0 {
		dao.DB().Where("id IN ?", missCommentIDs).
			Scopes(commonScopes...).
			Find(&linkedComments)
	}

	// Linked comments could be invisible but included in the list
	for _, c := range linkedComments {
		rCooked := dao.CookComment(&c)
		rCooked.Visible = false // Set invisible
		comments = append(comments, rCooked)
	}

	return comments
}
