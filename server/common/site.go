package common

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/gofiber/fiber/v2"
)

func UseSite(c *fiber.Ctx, siteName *string, destID *uint, destSiteAll *bool) {
	if destID != nil {
		*destID = c.Locals(config.CTX_KEY_ATK_SITE_ID).(uint)
	}
	if siteName != nil {
		*siteName = c.Locals(config.CTX_KEY_ATK_SITE_NAME).(string)
	}
	if destSiteAll != nil {
		*destSiteAll = c.Locals(config.CTX_KEY_ATK_SITE_ALL).(bool)
	}
}
