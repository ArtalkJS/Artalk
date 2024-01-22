package handler

import (
	"bytes"

	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ParamsTransferImport struct {
	artransfer.ImportParams
}

// @Id           ImportArtrans
// @Summary      Import Artrans
// @Description  Import data to Artalk
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Param        data  body  ParamsTransferImport  true  "The data to import"
// @Accept       json
// @Produce      html
// @Success      200  {string}  string
// @Router       /transfer/import  [post]
func TransferImport(app *core.App, router fiber.Router) {
	router.Post("/transfer/import", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var p ParamsTransferImport
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// TODO: temporary solution: output real-time log by html format body stream
		// consider using websocket or long polling in the future
		r := c.Response()

		r.Header.Add(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		r.SetStatusCode(fiber.StatusOK)

		buf := bytes.NewBufferString("")
		r.SetBodyStream(buf, -1)
		r.ImmediateHeaderFlush = true

		buf.Write([]byte(
			`<style>* { font-family: Menlo, Consolas, Monaco, monospace;word-wrap: break-word;white-space: pre-wrap;font-size: 13px; }</style>
		<script>function scroll() { if (!!document.body) { document.body.scrollTo(0, 999999999999); } }</script>`))

		artransfer.HttpOutput = func(continueRun bool, text string) {
			buf.Write([]byte(text))
			buf.Write([]byte("<script>scroll();</script>"))
		}

		p.Assumeyes = true
		artransfer.RunImportArtrans(app.Dao(), &p.ImportParams)

		return nil
	}))
}
