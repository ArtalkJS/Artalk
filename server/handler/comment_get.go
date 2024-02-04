package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseCommentGet struct {
	Comment      entity.CookedComment  `json:"comment"`       // The comment detail
	ReplyComment *entity.CookedComment `json:"reply_comment"` // The reply comment if exists (like reply)
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

		// Find comment by id
		comment := app.Dao().FindComment(uint(id))
		if comment.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}

		// Find linked comment by id
		var replyComment *entity.CookedComment
		if comment.Rid != 0 {
			rComment := app.Dao().FindComment(uint(comment.Rid))
			if !comment.IsEmpty() {
				rComment := app.Dao().CookComment(&rComment)
				rComment.Visible = false
				replyComment = &rComment
			}
		}

		cookedComment := app.Dao().CookComment(&comment)
		cookedComment = fetchIPRegionForComment(app, cookedComment)

		return common.RespData(c, ResponseCommentGet{
			Comment:      cookedComment,
			ReplyComment: replyComment,
		})
	})
}
