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

type ResponsePV struct {
	PV int `json:"pv"`
}

// @Summary      Page View
// @Description  Log and get the number of page views
// @Tags         PV
// @Param        page_key    formData  string  true   "the page key"
// @Param        page_title  formData  string  false  "the page title"
// @Param        site_name   formData  string  false  "the site name of your content scope"
// @Success      200  {object}  common.JSONResult{data=ResponsePV}
// @Router       /pv  [post]
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

		return common.RespData(c, ResponsePV{
			PV: page.PV,
		})
	})
}
