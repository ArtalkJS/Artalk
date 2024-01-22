package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           DeleteComment
// @Summary      Delete Comment
// @Description  Delete a specific comment
// @Tags         Comment
// @Param        id       path  int               true  "The comment ID you want to delete"
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /comments/{id}  [delete]
func CommentDelete(app *core.App, router fiber.Router) {
	router.Delete("/comments/:id", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		// find comment
		comment := app.Dao().FindComment(uint(id))
		if comment.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}

		// 删除主评论
		if err := app.Dao().DelComment(&comment); err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Comment")}))
		}

		// 删除子评论
		if err := app.Dao().DelCommentChildren(comment.ID); err != nil {
			return common.RespError(c, 500, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Sub-comment")}))
		}

		return common.RespSuccess(c)
	}))
}
