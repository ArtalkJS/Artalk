package http

import (
	"encoding/json"
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/artransfer"
	"github.com/labstack/echo/v4"
)

type ParamsAdminArtransfer struct {
	Type    string `mapstructure:"type" param:"required"`
	Payload string `mapstructure:"payload"`
}

func ActionAdminArtransfer(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminArtransfer
	if isOK, resp := ParamsDecode(c, ParamsAdminArtransfer{}, &p); !isOK {
		return resp
	}

	var payloadMap map[string]interface{}
	err := json.Unmarshal([]byte(p.Payload), &payloadMap)
	if err != nil {
		return RespError(c, "payload 解析错误", Map{
			"error": err,
		})
	}
	payload := []string{}
	for k, v := range payloadMap {
		payload = append(payload, k+":"+lib.ToString(v))
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	artransfer.Assumeyes = true
	artransfer.HttpOutput = func(continueRun bool, text string) {
		c.Response().Write([]byte(text))
		c.Response().Flush()
	}
	artransfer.RunByName(p.Type, payload)

	return nil
}
