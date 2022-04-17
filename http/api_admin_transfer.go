package http

import (
	"encoding/json"
	"net/http"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/artransfer"
	"github.com/labstack/echo/v4"
)

type ParamsAdminImport struct {
	Payload string `mapstructure:"payload"`
}

func ActionAdminImport(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsAdminImport
	if isOK, resp := ParamsDecode(c, ParamsAdminImport{}, &p); !isOK {
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

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	c.Response().Write([]byte(
		`<style>* { font-family: Menlo, Consolas, Monaco, monospace;word-wrap: break-word;white-space: pre-wrap;font-size: 13px; }</style>
		<script>function scroll() { document.body.scrollTo(0, 999999999999); }</script>`))
	c.Response().Flush()

	artransfer.Assumeyes = true
	artransfer.HttpOutput = func(continueRun bool, text string) {
		c.Response().Write([]byte(text))
		c.Response().Write([]byte("<script>scroll();</script>"))
		c.Response().Flush()
	}
	artransfer.RunImportArtrans(payload)

	return nil
}

func ActionAdminExport(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	jsonStr, err := artransfer.ExportArtransString()
	if err != nil {
		RespError(c, "导出错误", Map{
			"err": err,
		})
	}

	return RespData(c, Map{
		"data": jsonStr,
	})
}
