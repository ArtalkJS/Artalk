package handler

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminPageFetch struct {
	SiteName  string `json:"site_name"`  // The site name of your content scope
	GetStatus bool   `json:"get_status"` // If true, only get the status of the current task status
}

type ResponseAdminPageFetch struct {
	Data entity.CookedPage `json:"data"`
}

var allPageFetching = false
var allPageFetchDone = 0
var allPageFetchTotal = 0

// @Summary      Fetch Page Data
// @Description  Fetch the data of a specific page
// @Tags         Page
// @Security     ApiKeyAuth
// @Param        id       path  int                   true  "The page ID you want to fetch"
// @Param        options  body  ParamsAdminPageFetch  true  "The options"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseAdminPageFetch
// @Failure      400  {object}  Map{msg=string}
// @Failure      403  {object}  Map{msg=string}
// @Failure      404  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /pages/{id}/fetch  [post]
func AdminPageFetch(app *core.App, router fiber.Router) {
	router.Post("/pages/:id/fetch", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")

		var p ParamsAdminPageFetch
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// 状态获取
		if p.GetStatus {
			if allPageFetching {
				return common.RespData(c, common.Map{
					"msg":         i18n.T("{{done}} of {{total}} done", Map{"done": allPageFetchDone, "total": allPageFetchTotal}),
					"is_progress": true,
				})
			} else {
				return common.RespData(c, common.Map{
					"msg":         "",
					"is_progress": false,
				})
			}
		}

		// 更新全部站点
		// TODO separate the API `/pages/:id/fetch` and `/pages/fetch`
		if p.SiteName != "" {
			if allPageFetching {
				return common.RespError(c, 400, i18n.T("Task in progress, please wait a moment"))
			}

			// 异步执行
			go func() {
				allPageFetching = true
				allPageFetchDone = 0
				allPageFetchTotal = 0
				var pages []entity.Page
				db := app.Dao().DB().Model(&entity.Page{})
				if p.SiteName != config.ATK_SITE_ALL {
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
		}

		page := app.Dao().FindPageByID(uint(id))
		if page.IsEmpty() {
			return common.RespError(c, 404, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if !common.IsAdminHasSiteAccess(app, c, page.SiteName) {
			return common.RespError(c, 403, i18n.T("Access denied"))
		}

		if err := app.Dao().FetchPageFromURL(&page); err != nil {
			return common.RespError(c, 500, i18n.T("Page fetch failed")+": "+err.Error())
		}

		return common.RespData(c, ResponseAdminPageFetch{
			Data: app.Dao().CookPage(&page),
		})
	})
}
