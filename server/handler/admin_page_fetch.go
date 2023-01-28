package handler

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ParamsAdminPageFetch struct {
	ID       uint `form:"id"`
	SiteName string

	GetStatus bool `form:"get_status"`
}

var allPageFetching = false
var allPageFetchDone = 0
var allPageFetchTotal = 0

// @Summary      Page Data Fetch
// @Description  Fetch the data of a specific page
// @Tags         Page
// @Param        key            formData  string  true   "the page ID you want to fetch"
// @Param        site_name      formData  string  false  "the site name of your content scope"
// @Param        get_status     formData  bool    false  "which response data you want to receive"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/page-fetch  [post]
func AdminPageFetch(router fiber.Router) {
	router.Post("/page-fetch", func(c *fiber.Ctx) error {
		var p ParamsAdminPageFetch
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		common.UseSite(c, &p.SiteName, nil, nil)

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
		if p.SiteName != "" {
			if allPageFetching {
				return common.RespError(c, i18n.T("Task in progress, please wait a moment"))
			}

			// 异步执行
			go func() {
				allPageFetching = true
				allPageFetchDone = 0
				allPageFetchTotal = 0
				var pages []entity.Page
				db := db.DB().Model(&entity.Page{})
				if p.SiteName != config.ATK_SITE_ALL {
					db = db.Where(&entity.Page{SiteName: p.SiteName})
				}
				db.Find(&pages)

				allPageFetchTotal = len(pages)
				for _, p := range pages {
					if err := query.FetchPageFromURL(&p); err != nil {
						logrus.Error(c, "[api_admin_page_fetch] page fetch error: "+err.Error())
					} else {
						allPageFetchDone++
					}
				}
				allPageFetching = false
			}()

			return common.RespSuccess(c)
		}

		page := query.FindPageByID(p.ID)
		if page.IsEmpty() {
			return common.RespError(c, i18n.T("{{name}} not found", Map{"name": i18n.T("Page")}))
		}

		if !common.IsAdminHasSiteAccess(c, page.SiteName) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		if err := query.FetchPageFromURL(&page); err != nil {
			return common.RespError(c, i18n.T("Page fetch failed")+": "+err.Error())
		}

		return common.RespData(c, common.Map{
			"page": query.CookPage(&page),
		})
	})
}
