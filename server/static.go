package server

import "github.com/gofiber/fiber/v2"

type staticServer struct {
	app *fiber.App
}

func (s *staticServer) setupRoutes() {
	s.app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./views/dist/index.html", true)
	})
}

func newStaticServer() *staticServer {
	s := &staticServer{
		app: fiber.New(),
	}
	s.setupRoutes()
	return s
}
