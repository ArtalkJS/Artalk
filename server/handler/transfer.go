package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/gofiber/fiber/v2"
)

func AdminTransfer(app *core.App, router fiber.Router) {
	router.Post("/transfer/import", transferImport(app))
	router.Post("/transfer/upload", transferUpload(app))
	router.Get("/transfer/export", transferExport(app))
}
