package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPageFetchAll struct {
	SiteName string `json:"site_name" validate:"optional"` // If not empty, only fetch pages of this site
}

// @Id           FetchAllPages
// @Summary      Fetch All Pages Data
// @Description  Fetch the data of all pages
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        options  body  ParamsPageFetchAll  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/fetch  [post]
func PageFetchAll(app *core.App, router fiber.Router) {
	router.Post("/pages/fetch", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsPageFetchAll
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// If the task is in progress
		if allPageFetching {
			return common.RespError(c, 400, i18n.T("Task in progress, please wait a moment"))
		}

		// Start the async task
		go func() {
			allPageFetching = true
			allPageFetchDone = 0
			allPageFetchTotal = 0
			var pages []entity.Page
			db := app.Dao().DB().Model(&entity.Page{})
			if p.SiteName != "" {
				db = db.Where(&entity.Page{SiteName: p.SiteName})
			}
			db.Find(&pages)

			allPageFetchTotal = len(pages)
			for _, p := range pages {
				if err := app.Dao().FetchPageFromURL(&p); err != nil {
					log.Error(c, "[api_admin_page_fetch] page fetch error: "+err.Error())
				} else {
					allPageFetchDone++
				}
			}
			allPageFetching = false
		}()

		return common.RespSuccess(c)
	}))
}
