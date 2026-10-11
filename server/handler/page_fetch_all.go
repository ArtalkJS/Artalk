package handler

import (
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/server/common"
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

		// Reserve the task before starting the goroutine, so simultaneous requests
		// cannot launch duplicate batches. Status readers share the same lock.
		if !allPageFetch.start() {
			return common.RespError(c, 400, i18n.T("Task in progress, please wait a moment"))
		}

		// Start the async task
		go func() {
			defer allPageFetch.finish()
			var pages []entity.Page
			db := app.Dao().DB().Model(&entity.Page{})
			if p.SiteName != "" {
				db = db.Where(&entity.Page{SiteName: p.SiteName})
			}
			if err := db.Find(&pages).Error; err != nil {
				log.Error("[api_admin_page_fetch] failed to list pages: ", err)
				return
			}

			done := 0
			allPageFetch.progress(done, len(pages))
			for _, p := range pages {
				if err := app.Dao().FetchPageFromURL(&p, app.Conf().TrustedDomains...); err != nil {
					log.Error("[api_admin_page_fetch] page fetch error: ", err)
				} else {
					done++
					allPageFetch.progress(done, len(pages))
				}
			}
		}()

		return common.RespSuccess(c)
	}))
}
