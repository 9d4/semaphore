package server

import "github.com/gofiber/fiber/v2"

var app fiber.App

func Start() error {
	app = *fiber.New()

	setupRoutes(&app)

	return app.Listen("0.0.0.0:3500")
}
