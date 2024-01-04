package handler

import (
	"bytes"
	"io"
	"os"
	"slices"

	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

func AdminTransfer(app *core.App, router fiber.Router) {
	router.Post("/transfer/import", transferImport(app))
	router.Post("/transfer/upload", transferUpload(app))
	router.Get("/transfer/export", transferExport(app))
}

type ParamsAdminImport struct {
	artransfer.ImportParams
}

// @Summary      Import Artrans
// @Description  Import data to Artalk
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Param        data  body  ParamsAdminImport  true  "The data to import"
// @Accept       json
// @Produce      html
// @Success      200  {string}  string
// @Router       /transfer/import  [post]
func transferImport(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var p ParamsAdminImport
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// If not super admin, force to fill target site name and check permission
		if !common.GetIsSuperAdmin(app, c) {
			if p.TargetSiteName == "" {
				return common.RespError(c, 400, "Please fill in the target site name")
			}

			user := common.GetUserByReq(app, c)
			if !slices.Contains(app.Dao().CookUser(&user).SiteNames, p.TargetSiteName) {
				return common.RespError(c, 400, "Destination site name are prohibited since no permission")
			}
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
	}
}

type ResponseExport struct {
	// The exported data which is a JSON string
	Artrans string `json:"artrans"`
}

// @Summary      Export Artrans
// @Description  Export data from Artalk
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseExport
// @Failure      500  {object}  Map{msg=string}
// @Router       /transfer/export  [get]
func transferExport(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var siteNameScope []string

		// If not super admin, only export sites that have permission
		if !common.GetIsSuperAdmin(app, c) {
			u := common.GetUserByReq(app, c)
			siteNameScope = app.Dao().CookUser(&u).SiteNames
		}

		jsonStr, err := artransfer.RunExportArtrans(app.Dao(), &artransfer.ExportParams{
			SiteNameScope: siteNameScope,
		})
		if err != nil {
			common.RespError(c, 500, i18n.T("Export error"), common.Map{
				"err": err,
			})
		}

		return common.RespData(c, ResponseExport{
			Artrans: jsonStr,
		})
	}
}

type ResponseImportUpload struct {
	// The uploaded file name which can be used to import
	Filename string `json:"filename"`
}

// @Summary      Upload Artrans
// @Description  Upload a file to prepare to import
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Param        file  formData  file  true  "Upload file in preparation for import task"
// @Accept       mpfd
// @Produce      json
// @Success      200  {object}  ResponseImportUpload{filename=string}
// @Failure      500  {object}  Map{msg=string}
// @Router       /transfer/upload  [post]
func transferUpload(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get file from FormData
		file, err := c.FormFile("file")
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File read failed")
		}

		// Open file
		src, err := file.Open()
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File open failed")
		}
		defer src.Close()

		// Read file to buffer
		buf, err := io.ReadAll(src)
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "File read failed")
		}

		// Create temp file
		tmpFile, err := os.CreateTemp("", "artalk-import-file-")
		if err != nil {
			log.Error(err)
			return common.RespError(c, 500, "tmp file creation failed")
		}

		// Write buffer to temp file
		tmpFile.Write(buf)

		// Return filename
		return common.RespData(c, ResponseImportUpload{
			Filename: tmpFile.Name(),
		})
	}
}
