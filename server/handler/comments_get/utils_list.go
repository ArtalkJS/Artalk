package comments_get

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

// Find all child comments (for nested mode)
func FindChildComments(dao *dao.Dao, user entity.User, comments []entity.CookedComment) []entity.CookedComment {
	for _, parent := range comments { // TODO: Consider add a feature, read more children, pagination for children comment
		children := dao.FindCommentChildren(parent.ID, NoPendingChecker(user))
		comments = append(comments, dao.CookAllComments(children)...)
	}

	return comments
}

// Find all linked comments (for flat mode)
func FindLinkedComments(app *dao.Dao, comments []entity.CookedComment) []entity.CookedComment {
	// Find linked comments which `id` is the same as `rid`
	// Linked comments could be invisible but included in the list
	for _, comment := range comments {
		// If comment is root comment, skip
		if comment.Rid == 0 {
			continue
		}

		// If comment is already in the list, skip
		if entity.ContainsCookedComment(comments, comment.Rid) {
			continue
		}

		// Get linked comment
		rComment := app.FindComment(comment.Rid)
		if rComment.IsEmpty() {
			continue
		}

		// Set invisible
		rCooked := app.CookComment(&rComment)
		rCooked.Visible = false

		comments = append(comments, rCooked)
	}

	return comments
}

// Find the IP region of each comment
func FindIPRegionForComments(app *core.App, comments []entity.CookedComment) []entity.CookedComment {
	if !app.Conf().IPRegion.Enabled {
		return comments
	}

	ipRegionService, err := core.AppService[*core.IPRegionService](app)
	if err == nil {
		for _, c := range comments {
			c.IPRegion = ipRegionService.Query(c.IP)
		}
	} else {
		log.Error("[IPRegionService] err: ", err)
	}

	return comments
}
