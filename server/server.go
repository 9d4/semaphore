package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var app fiber.App

func Start() error {
	app = *fiber.New()

	setupRoutes(&app)

	return app.Listen(viper.GetString("addr"))
}
