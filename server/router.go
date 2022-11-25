package server

import "github.com/gofiber/fiber/v2"

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		c.WriteString("Merhaba 👋🏼")
		return nil
	})
}
