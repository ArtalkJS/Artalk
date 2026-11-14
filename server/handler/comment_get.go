package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseCommentGet struct {
	Comment      entity.CommentDTO  `json:"comment"`       // The comment detail
	ReplyComment *entity.CommentDTO `json:"reply_comment"` // The reply comment if exists (like reply)
}

// @Id           GetComment
// @Summary      Get a comment
// @Description  Get the detail of a comment by comment id
// @Tags         Comment
// @Param        id       path  int  true  "The comment ID you want to get"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseCommentGet
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments/{id}  [get]
func CommentGet(app *core.App, router fiber.Router) {
	router.Get("/comments/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		user, err := common.GetUserByReq(app, c)
		if err != nil {
			user = entity.User{}
		}
		guestHash := common.GuestTokenHash(c)
		canRead := func(comment entity.Comment) bool {
			return !comment.IsEmpty() && (!comment.IsPending || user.IsAdmin ||
				(user.ID != 0 && user.ID == comment.UserID) ||
				(guestHash != "" && guestHash == comment.GuestTokenHash))
		}

		// Find comment by id
		comment := app.Dao().FindComment(uint(id))
		if !canRead(comment) {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}

		// Find linked comment by id
		var replyComment *entity.CommentDTO
		if comment.Rid != 0 {
			rComment := app.Dao().FindComment(uint(comment.Rid))
			if canRead(rComment) {
				rComment := app.Dao().BuildCommentDTO(&rComment)
				rComment.Visible = false
				replyComment = &rComment
			}
		}

		commentDTO := app.Dao().BuildCommentDTO(&comment)
		commentDTO = fetchIPRegionForComment(app, commentDTO)

		return common.RespData(c, ResponseCommentGet{
			Comment:      commentDTO,
			ReplyComment: replyComment,
		})
	})
}
