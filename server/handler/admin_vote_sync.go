package handler

import (
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminVoteSync struct {
}

// @Summary      Vote Sync
// @Description  Sync the number of votes in the `comments` or `pages` data tables to keep them the same as the `votes` table
// @Tags         Vote
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/vote-sync  [post]
func AdminVoteSync(router fiber.Router) {
	router.Post("/vote-sync", func(c *fiber.Ctx) error {
		var p ParamsAdminVoteSync
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		query.VoteSync()

		return common.RespSuccess(c)
	})
}
