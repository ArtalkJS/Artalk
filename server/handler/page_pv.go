package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPagePV struct {
	PageKey   string `json:"page_key" validate:"required"`   // The page key
	PageTitle string `json:"page_title" validate:"optional"` // The page title
	SiteName  string `json:"site_name" validate:"optional"`  // The site name of your content scope
}

type ResponsePagePV struct {
	PV int `json:"pv"`
}

// @Id           LogPv
// @Summary      Increase Page Views (PV)
// @Description  Increase and get the number of page views
// @Tags         Page
// @Param        page  body  ParamsPagePV  true  "The page to record pv"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponsePagePV
// @Router       /pages/pv  [post]
func PagePV(app *core.App, router fiber.Router) {
	router.Post("/pages/pv", func(c *fiber.Ctx) error {
		var p ParamsPagePV
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// find page
		page := app.Dao().FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

		// ip := c.RealIP()
		// ua := c.Request().UserAgent()

		page.PV++
		app.Dao().UpdatePage(&page)

		return common.RespData(c, ResponsePagePV{
			PV: page.PV,
		})
	})
}
