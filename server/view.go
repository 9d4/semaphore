package server

import "github.com/gofiber/fiber/v2"

// viewApp will always serve Html file of SPA frontend.
func viewApp() *fiber.App {
	app := fiber.New()
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./views/dist/index.html", true)
	})
	return app
}
