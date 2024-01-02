package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

func AdminTransfer(app *core.App, router fiber.Router) {
	router.Post("/artransfer/import", adminImport(app))
	router.Post("/artransfer/upload", adminImportUpload(app))
	router.Get("/artransfer/export", adminExport(app))
}

type ParamsAdminImport struct {
	Payload string `json:"payload"` // The transfer importer payload
}

// @Summary      Import Artrans
// @Description  Import data to Artalk
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Param        data  body  ParamsAdminImport  true  "The data to import"
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.JSONResult
// @Router       /artransfer/import  [post]
func adminImport(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var p ParamsAdminImport
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		var payloadMapRaw map[string]interface{}
		err := json.Unmarshal([]byte(p.Payload), &payloadMapRaw)
		if err != nil {
			return common.RespError(c, "Payload parsing error", common.Map{
				"error": err,
			})
		}

		payloadMap := map[string]string{}
		for k, v := range payloadMapRaw {
			payloadMap[k] = utils.ToString(v) // convert all value to string
		}

		payloadArr := []string{}
		for k, v := range payloadMap {
			payloadArr = append(payloadArr, k+":"+v)
		}

		if !common.GetIsSuperAdmin(app, c) {
			user := common.GetUserByReq(app, c)
			if sitName, isExist := payloadMap["t_name"]; isExist {
				if !utils.ContainsStr(app.Dao().CookUser(&user).SiteNames, sitName) {
					return common.RespError(c, "Destination site name of prohibited import")
				}
			} else {
				return common.RespError(c, "Please fill in the target site name")
			}
		}

		// TODO bcz 懒，先整这个缓冲输出，以后改成高级点的
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

		params := artransfer.ArrToImportParams(payloadArr)
		params.Assumeyes = true
		artransfer.RunImportArtrans(app.Dao(), params)

		return nil
	}
}

// @Summary      Upload Artrans
// @Description  Upload a file to prepare to import
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Param        file  formData  file  true  "Upload file in preparation for import task"
// @Accept       mpfd
// @Produce      json
// @Success      200  {object}  common.JSONResult{data=object{filename=string}}
// @Router       /artransfer/upload  [post]
func adminImportUpload(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// 获取 Form
		file, err := c.FormFile("file")
		if err != nil {
			log.Error(err)
			return common.RespError(c, "File read failed")
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			log.Error(err)
			return common.RespError(c, "File open failed")
		}
		defer src.Close()

		// 读取文件
		buf, err := io.ReadAll(src)
		if err != nil {
			log.Error(err)
			return common.RespError(c, "File read failed")
		}

		tmpFile, err := os.CreateTemp("", "artalk-import-file-")
		if err != nil {
			log.Error(err)
			return common.RespError(c, "tmp file creation failed")
		}

		tmpFile.Write(buf)

		return common.RespData(c, common.Map{
			"filename": tmpFile.Name(),
		})
	}
}

// @Summary      Export Artrans
// @Description  Export data from Artalk
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  common.JSONResult{data=object{data=string}}
// @Router       /artransfer/export  [get]
func adminExport(app *core.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var siteNameScope []string
		if !common.GetIsSuperAdmin(app, c) {
			// 仅导出限定范围内的站点
			u := common.GetUserByReq(app, c)
			siteNameScope = app.Dao().CookUser(&u).SiteNames
		}

		jsonStr, err := artransfer.RunExportArtrans(app.Dao(), &artransfer.ExportParams{
			SiteNameScope: siteNameScope,
		})
		if err != nil {
			common.RespError(c, i18n.T("Export error"), common.Map{
				"err": err,
			})
		}

		return common.RespData(c, common.Map{
			"data": jsonStr,
		})
	}
}
