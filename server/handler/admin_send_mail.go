package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSendMail struct {
	Subject string `json:"subject" validate:"required"` // The subject of email
	Body    string `json:"body" validate:"required"`    // The body of email
	ToAddr  string `json:"to_addr" validate:"required"` // The email address of the receiver
}

// @Summary      Send Email
// @Description  Send an email to test the email sender
// @Tags         System
// @Security     ApiKeyAuth
// @Param        email  body  ParamsAdminSendMail  true  "The email data"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Map{}
// @Failure      403  {object}  Map{msg=string}
// @Failure      500  {object}  Map{}
// @Router       /send_email  [post]
func AdminSendMail(app *core.App, router fiber.Router) {
	router.Post("/send_email", common.AdminGuard(app, func(c *fiber.Ctx) error {
		if ok, resp := common.AdminRequired(app, c); !ok {
			return resp
		}

		var p ParamsAdminSendMail
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(app, c) {
			return common.RespError(c, 403, i18n.T("Access denied"))
		}

		emailService, err := core.AppService[*core.EmailService](app)
		if err != nil {
			log.Error("[EmailService] err: ", err)
			return common.RespError(c, 500, "EmailService err: "+err.Error())
		}
		emailService.AsyncSendTo(p.Subject, p.Body, p.ToAddr)

		return common.RespSuccess(c)
	}))
}
