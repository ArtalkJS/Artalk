package handler

import (
	"github.com/ArtalkJS/Artalk/internal/email"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSendMail struct {
	Subject string `form:"subject" validate:"required"`
	Body    string `form:"body" validate:"required"`
	ToAddr  string `form:"to_addr" validate:"required"`
}

// @Summary      Email Send
// @Description  Send an email to test the email sender
// @Tags         System
// @Param        subject        formData  string  true  "the subject of email"
// @Param        body           formData  string  true  "the body of email"
// @Param        to_addr        formData  string  true  "the email address of the receiver"
// @Security     ApiKeyAuth
// @Success      200  {object}  common.JSONResult
// @Router       /admin/send-mail  [post]
func AdminSendMail(router fiber.Router) {
	router.Post("/send-mail", func(c *fiber.Ctx) error {
		var p ParamsAdminSendMail
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, i18n.T("Access denied"))
		}

		email.AsyncSendTo(p.Subject, p.Body, p.ToAddr)

		return common.RespSuccess(c)
	})
}
