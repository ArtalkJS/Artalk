package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseAdminPageFetch struct {
	entity.CookedPage
}

var allPageFetching = false
var allPageFetchDone = 0
var allPageFetchTotal = 0

// @Summary      Fetch Page Data
// @Description  Fetch the data of a specific page
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        id       path  int                   true  "The page ID you want to fetch"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminPageFetch
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/{id}/fetch  [post]
func AdminPageFetch(app *core.App, router fiber.Router) {
	router.Post("/pages/:id/fetch", common.AdminGuard(app, func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		page := app.Dao().FindPageByID(uint(id))
		if page.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if err := app.Dao().FetchPageFromURL(&page); err != nil {
			return common.RespError(c, 500, i18n.T("Page fetch failed")+": "+err.Error())
		}

		return common.RespData(c, ResponseAdminPageFetch{
			CookedPage: app.Dao().CookPage(&page),
		})
	}))
}

type ParamsAdminFetchAllPages struct {
	SiteName  string `json:"site_name"`  // If not empty, only fetch pages of this site
	GetStatus bool   `json:"get_status"` // If true, only get the status of the current task status
}

// @Summary      Fetch All Pages Data
// @Description  Fetch the data of all pages
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        options  body  ParamsAdminFetchAllPages  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/fetch  [post]
func AdminPagesFetchAll(app *core.App, router fiber.Router) {
	router.Post("/pages/fetch", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsAdminFetchAllPages
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// If the task is in progress
		if allPageFetching {
			// If user want to get the status
			if p.GetStatus {

			} else {
				return common.RespError(c, 400, i18n.T("Task in progress, please wait a moment"))
			}
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

type ResponseAdminPageFetchAllStatus struct {
	Msg        string `json:"msg"`         // The message of the task status
	IsProgress bool   `json:"is_progress"` // If the task is in progress
	Done       int    `json:"done"`        // The number of pages that have been fetched
	Total      int    `json:"total"`       // The total number of pages
}

// @Summary      Get All Pages Fetch Status
// @Description  Get the status of the task of fetching all pages
// @Tags         Page
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseAdminPageFetchAllStatus
// @Router       /pages/fetch/status  [get]
func AdminPageFetchAllStatus(app *core.App, router fiber.Router) {
	router.Get("/pages/fetch/status", common.AdminGuard(app, func(c *fiber.Ctx) error {
		return common.RespData(c, ResponseAdminPageFetchAllStatus{
			Msg:        i18n.T("{{done}} of {{total}} done", Map{"done": allPageFetchDone, "total": allPageFetchTotal}),
			IsProgress: allPageFetching,
			Done:       allPageFetchDone,
			Total:      allPageFetchTotal,
		})
	}))
}
