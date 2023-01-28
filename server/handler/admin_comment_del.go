package handler

import (
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsCommentDel struct {
	ID uint `form:"id" validate:"required"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

// @Summary      Comment Delete
// @Description  Delete a specific comment
// @Tags         Comment
// @Param        id             formData  int     true  "the comment ID you want to delete"
// @Param        site_name      formData  string  false "the site name of your content scope"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/comment-del  [post]
func AdminCommentDel(router fiber.Router) {
	router.Post("/comment-del", func(c *fiber.Ctx) error {
		var p ParamsCommentDel
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		// find comment
		comment := query.FindComment(p.ID)
		if comment.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Comment")}))
		}

		if !common.IsAdminHasSiteAccess(c, comment.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		// 删除主评论
		if err := query.DelComment(&comment); err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Comment")}))
		}

		// 删除子评论
		if err := query.DelCommentChildren(comment.ID); err != nil {
			return common.RespError(c, i18n.T("{{name}} deletion failed", Map{"name": i18n.T("Sub-comment")}))
		}

		return common.RespSuccess(c)
	})
}
