package http

import (
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/labstack/echo/v4"
)

type ParamsSendMail struct {
	Subject string `mapstructure:"subject" param:"required"`
	Body    string `mapstructure:"body" param:"required"`
	ToAddr  string `mapstructure:"to_addr" param:"required"`
}

func ActionManagerSendMail(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsSendMail
	if isOK, resp := ParamsDecode(c, ParamsSendMail{}, &p); !isOK {
		return resp
	}

	email.SendTo(p.Subject, p.Body, p.ToAddr)

	return RespSuccess(c)
}
