package handler

import (
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPV struct {
	PageKey   string `form:"page_key" validate:"required"`
	PageTitle string `form:"page_title"`

	SiteName string
	SiteID   uint
	SiteAll  bool
}

// POST /api/pv
func PV(router fiber.Router) {
	router.Post("/pv", func(c *fiber.Ctx) error {
		var p ParamsPV
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// use site
		common.UseSite(c, &p.SiteName, &p.SiteID, &p.SiteAll)

		// find page
		page := query.FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

		// ip := c.RealIP()
		// ua := c.Request().UserAgent()

		page.PV++
		query.UpdatePage(&page)

		return common.RespData(c, common.Map{
			"pv": page.PV,
		})
	})
}
