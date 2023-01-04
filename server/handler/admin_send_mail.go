package handler

import (
	"github.com/ArtalkJS/ArtalkGo/internal/email"
	"github.com/ArtalkJS/ArtalkGo/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsAdminSendMail struct {
	Subject string `form:"subject" validate:"required"`
	Body    string `form:"body" validate:"required"`
	ToAddr  string `form:"to_addr" validate:"required"`
}

// POST /api/admin/send-mail
func AdminSendMail(router fiber.Router) {
	router.Post("/send-mail", func(c *fiber.Ctx) error {
		var p ParamsAdminSendMail
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		if !common.GetIsSuperAdmin(c) {
			return common.RespError(c, "无权访问")
		}

		email.AsyncSendTo(p.Subject, p.Body, p.ToAddr)

		return common.RespSuccess(c)
	})
}
