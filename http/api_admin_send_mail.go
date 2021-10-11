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

func ActionAdminSendMail(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminSendMail
	if isOK, resp := ParamsDecode(c, ParamsAdminSendMail{}, &p); !isOK {
		return resp
	}

	email.AsyncSendTo(p.Subject, p.Body, p.ToAddr)

	return RespSuccess(c)
}
