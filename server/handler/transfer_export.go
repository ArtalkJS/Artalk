package handler

import (
	"github.com/ArtalkJS/Artalk/internal/artransfer"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/server/common"
	"github.com/gofiber/fiber/v2"
)

type ResponseTransferExport struct {
	// The exported data which is a JSON string
	Artrans string `json:"artrans"`
}

// @Id           ExportArtrans
// @Summary      Export Artrans
// @Description  Export data from Artalk
// @Tags         Transfer
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  ResponseTransferExport
// @Failure      500  {object}  Map{msg=string}
// @Router       /transfer/export  [get]
func TransferExport(app *core.App, router fiber.Router) {
	router.Get("/transfer/export", common.AdminGuard(app, func(c *fiber.Ctx) error {
		var siteNameScope []string

		jsonStr, err := artransfer.RunExportArtrans(app.Dao(), &artransfer.ExportParams{
			SiteNameScope: siteNameScope,
		})
		if err != nil {
			common.RespError(c, 500, i18n.T("Export error"), common.Map{
				"err": err,
			})
		}

		return common.RespData(c, ResponseTransferExport{
			Artrans: jsonStr,
		})
	}))
}
