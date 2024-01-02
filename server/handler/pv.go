package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsPV struct {
	PageKey   string `json:"page_key" validate:"required"` // The page key
	PageTitle string `json:"page_title"`                   // The page title
	SiteName  string `json:"site_name"`                    // The site name of your content scope
}

type ResponsePV struct {
	PV int `json:"pv"`
}

// @Summary      Record PV
// @Description  Log and get the number of page views
// @Tags         Page
// @Param        page  body  ParamsPV  true  "The page to record pv"
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.JSONResult{data=ResponsePV}
// @Router       /pages/pv  [post]
func PV(app *core.App, router fiber.Router) {
	router.Post("/pages/pv", func(c *fiber.Ctx) error {
		var p ParamsPV
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// find page
		page := app.Dao().FindCreatePage(p.PageKey, p.PageTitle, p.SiteName)

		// ip := c.RealIP()
		// ua := c.Request().UserAgent()

		page.PV++
		app.Dao().UpdatePage(&page)

		return common.RespData(c, ResponsePV{
			PV: page.PV,
		})
	})
}
