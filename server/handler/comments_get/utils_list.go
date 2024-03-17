package comments_get

import (
	"slices"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

// Find all child comments (for nested mode)
func FlattenChildComments(dao *dao.Dao, user entity.User, comments []*entity.Comment) []*entity.Comment {
	flatten := make([]*entity.Comment, 0)
	queue := make([]*entity.Comment, len(comments))
	copy(queue, comments)

	for len(queue) > 0 {
		c := queue[0]     // get the first element
		queue = queue[1:] // dequeue

		if !NoPendingChecker(user)(c) {
			continue
		}

		flatten = append(flatten, c)

		// add children to the end of the queue
		queue = append(queue, c.Children...)
	}

	return flatten
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
		if slices.ContainsFunc(comments, func(c entity.CookedComment) bool {
			return c.ID == comment.Rid
		}) {
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
		for i, c := range comments {
			comments[i].IPRegion = ipRegionService.Query(c.IP)
		}
	} else {
		log.Error("[IPRegionService] err: ", err)
	}

	return comments
}
