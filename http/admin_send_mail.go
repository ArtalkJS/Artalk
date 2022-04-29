package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/labstack/echo/v4"
)

type ParamsAdminSendMail struct {
	Subject string `mapstructure:"subject" param:"required"`
	Body    string `mapstructure:"body" param:"required"`
	ToAddr  string `mapstructure:"to_addr" param:"required"`
}

func (a *action) AdminSendMail(c echo.Context) error {
	var p ParamsAdminSendMail
	if isOK, resp := ParamsDecode(c, &p); !isOK {
		return resp
	}

	if !GetIsSuperAdmin(c) {
		return RespError(c, "无权访问")
	}

	email.AsyncSendTo(p.Subject, p.Body, p.ToAddr)

	return RespSuccess(c)
}
