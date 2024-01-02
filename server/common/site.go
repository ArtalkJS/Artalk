package common

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/gofiber/fiber/v2"
)

type SiteInfoByRequest struct {
	ID   uint
	Name string
	All  bool
}

func GetSiteInfo(c *fiber.Ctx) SiteInfoByRequest {
	return SiteInfoByRequest{
		ID:   c.Locals(config.CTX_KEY_ATK_SITE_ID).(uint),
		Name: c.Locals(config.CTX_KEY_ATK_SITE_NAME).(string),
		All:  c.Locals(config.CTX_KEY_ATK_SITE_ALL).(bool),
	}
}
