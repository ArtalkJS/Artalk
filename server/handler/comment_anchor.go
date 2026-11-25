package handler

import (
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCommentAnchorUpdate struct {
	Resolved *bool             `json:"resolved" validate:"optional"` // The annotation resolved status
	Anchor   *entity.AnchorDTO `json:"anchor" validate:"optional"`   // The annotation anchor
}

type ResponseCommentAnchorUpdate struct {
	entity.CommentDTO
}

// @Id           UpdateCommentAnchor
// @Summary      Update Comment Anchor
// @Description  Update the annotation anchor or resolved status of a comment
// @Tags         Comment
// @Param        id      path  int                        true  "The comment ID you want to update"
// @Param        anchor  body  ParamsCommentAnchorUpdate  true  "The anchor update"
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseCommentAnchorUpdate
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments/{id}/anchor  [put]
func CommentAnchor(app *core.App, router fiber.Router) {
	router.Put("/comments/:id/anchor", func(c *fiber.Ctx) error {
		var p ParamsCommentAnchorUpdate
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if p.Resolved == nil && p.Anchor == nil {
			return common.RespError(c, 400, "Nothing to update")
		}

		id, _ := c.ParamsInt("id")
		comment := app.Dao().FindComment(uint(id))
		if comment.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}
		if !comment.IsAnnotationThread || comment.Rid != 0 {
			return common.RespError(c, 400, "Comment is not an annotation")
		}
		if !canMutateAnnotationComment(app, c, &comment) {
			return common.RespError(c, 403, i18n.T("Admin access required"), Map{"need_login": true})
		}

		if p.Resolved != nil {
			comment.AnchorResolved = *p.Resolved
		}
		if p.Anchor != nil {
			if err := validateAnchor(p.Anchor); err != nil {
				return common.RespError(c, 400, err.Error())
			}
			comment.AnchorExact = strings.TrimSpace(p.Anchor.Quote.Exact)
			comment.AnchorPrefix = p.Anchor.Quote.Prefix
			comment.AnchorSuffix = p.Anchor.Quote.Suffix
			comment.AnchorSection = p.Anchor.Section
			comment.AnchorStart = -1
			comment.AnchorEnd = -1
			if p.Anchor.Position != nil {
				comment.AnchorStart = p.Anchor.Position.Start
				comment.AnchorEnd = p.Anchor.Position.End
			}
		}

		if err := app.Dao().UpdateComment(&comment); err != nil {
			log.Error("[CommentAnchor] update error: ", err)
			return common.RespError(c, 500, i18n.T("{{name}} save failed", Map{"name": i18n.T("Comment")}))
		}

		dto := app.Dao().BuildCommentDTO(&comment)
		return common.RespData(c, ResponseCommentAnchorUpdate{CommentDTO: dto})
	})
}

// canMutateAnnotationComment allows admins to mutate any comment, and authors
// to mutate comments that belong to an annotation thread.
func canMutateAnnotationComment(app *core.App, c *fiber.Ctx, comment *entity.Comment) bool {
	if common.CheckIsAdminReq(app, c) {
		return true
	}
	if !comment.IsAnnotationThread {
		return false
	}
	user, _ := common.GetUserByReq(app, c)
	if !user.IsEmpty() && comment.UserID == user.ID {
		return true
	}
	hash := common.GuestTokenHash(c)
	return hash != "" && hash == comment.GuestTokenHash
}
