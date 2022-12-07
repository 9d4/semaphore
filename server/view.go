package server

import "github.com/gofiber/fiber/v2"

type viewServer struct {
	app *fiber.App
}

func (s *viewServer) setupRoutes() {
	s.app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./views/dist/index.html", true)
	})
}

func newViewServer() *viewServer {
	s := &viewServer{
		app: fiber.New(),
	}
	s.setupRoutes()
	return s
}
