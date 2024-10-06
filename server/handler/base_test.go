package handler_test

import (
	"github.com/artalkjs/artalk/v2/test"
	"github.com/gofiber/fiber/v2"
)

func NewApiTestApp() (*test.TestApp, *fiber.App) {
	app, _ := test.NewTestApp()
	fiberApp := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})
	return app, fiberApp
}
