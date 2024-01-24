package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseCaptchaStatus struct {
	IsPass bool `json:"is_pass"`
}

// @Id           GetCaptchaStatus
// @Summary      Get Captcha Status
// @Description  Get the status of the user's captcha verification
// @Tags         Captcha
// @Produce      json
// @Success      200  {object}  ResponseCaptchaStatus
// @Router       /captcha/status  [get]
func CaptchaStatus(app *core.App, router fiber.Router) {
	router.Get("/captcha/status", func(c *fiber.Ctx) error {
		limiter, err := common.GetLimiter(c)
		if err != nil {
			return err
		}

		return common.RespData(c, ResponseCaptchaStatus{
			IsPass: limiter.IsPass(c.IP()),
		})
	})
}
