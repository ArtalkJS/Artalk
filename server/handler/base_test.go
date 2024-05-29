package handler_test

import (
	"github.com/ArtalkJS/Artalk/test"
	"github.com/gofiber/fiber/v2"
)

func NewApiTestApp() (*test.TestApp, *fiber.App) {
	app, _ := test.NewTestApp()
	fiberApp := fiber.New()
	return app, fiberApp
}
