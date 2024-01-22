package handler

import (
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/gofiber/fiber/v2"
)

func Transfer(app *core.App, router fiber.Router) {
	TransferImport(app, router)
	TransferUpload(app, router)
	TransferExport(app, router)
}
