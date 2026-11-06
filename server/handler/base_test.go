package handler_test

import (
	"github.com/artalkjs/artalk/v2/test"
	"github.com/gofiber/fiber/v2"
)

func NewApiTestApp() (*test.TestApp, *fiber.App) {
	app, _ := test.NewTestApp()
	fiberApp := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
		Immutable:   true, // Match production: request strings outlive the response in background jobs.
	})
	return app, fiberApp
}
