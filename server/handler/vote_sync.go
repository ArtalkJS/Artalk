package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

// @Id           SyncVotes
// @Summary      Sync Vote Data
// @Description  Sync the number of votes in the `comments` or `pages` data tables to keep them the same as the `votes` table
// @Tags         Vote
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Router       /votes/sync  [post]
func VoteSync(app *core.App, router fiber.Router) {
	router.Post("/votes/sync", common.AdminGuard(app, func(c *fiber.Ctx) error {
		app.Dao().VoteSync()

		return common.RespSuccess(c)
	}))
}
