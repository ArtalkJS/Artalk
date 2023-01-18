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

// POST /api/admin/comment-del
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
			return common.RespError(c, i18n.T("Comment not found"))
		}

		if !common.IsAdminHasSiteAccess(c, comment.SiteName) {
			return common.RespError(c, i18n.T("No access"))
		}

		// 删除主评论
		if err := query.DelComment(&comment); err != nil {
			return common.RespError(c, i18n.T("Failed to delete comment"))
		}

		// 删除子评论
		if err := query.DelCommentChildren(comment.ID); err != nil {
			return common.RespError(c, i18n.T("Sub-comment deletion failed"))
		}

		return common.RespSuccess(c)
	})
}
