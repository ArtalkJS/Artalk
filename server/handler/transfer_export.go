package handler

import (
	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

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
	return common.AdminGuard(app, func(c *fiber.Ctx) error {
		var siteNameScope []string

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
	})
}
