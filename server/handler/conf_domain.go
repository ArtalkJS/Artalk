package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/ArtalkJS/Artalk/server/middleware"
	"github.com/gofiber/fiber/v2"
)

type ResponseConfDomain struct {
	Origin    string `json:"origin"`     // The origin of the domain
	IsTrusted bool   `json:"is_trusted"` // Is the domain trusted
}

// @Id            GetDomain
// @Summary       Get Domain Info
// @Description   Get Domain Info
// @Tags          System
// @Produce       json
// @Param         url query string false "Domain URL"
// @Success       200  {object}  ResponseConfDomain
// @Router        /conf/domain  [get]
func ConfDomain(app *core.App, router fiber.Router) {
	router.Get("/conf/domain", func(c *fiber.Ctx) error {
		domainURL := c.Query("url")
		if domainURL == "" {
			domainURL = c.Get("Origin")
		}
		trusted, origin, _ := middleware.CheckURLTrusted(app, domainURL)
		return common.RespData(c, ResponseConfDomain{
			IsTrusted: trusted,
			Origin:    origin,
		})
	})
}
