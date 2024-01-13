package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/limiter"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseCaptchaStatus struct {
	IsPass bool `json:"is_pass"`
}

// @Summary      Get Captcha Status
// @Description  Get the status of the user's captcha verification
// @Tags         Captcha
// @Produce      json
// @Success      200  {object}  ResponseCaptchaStatus
// @Router       /captcha/status  [get]
func captchaStatus(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		limiter, err := common.GetLimiter[limiter.Limiter](c)
		if limiter == nil {
			return err
		}

		return common.RespData(c, ResponseCaptchaStatus{
			IsPass: limiter.IsPass(c.IP()),
		})
	}
}
