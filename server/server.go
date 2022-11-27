package server

import (
	"github.com/9d4/semaphore/store"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type server struct {
	app   *fiber.App
	db    *gorm.DB
	store store.Store
	v     *viper.Viper
}

func (s *server) setupRoutes() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		c.WriteString("Merhaba 👋🏼")
		return nil
	})

	apiSrv := newApiServer(s.db, s.store)
	s.app.Mount("/api", apiSrv.app)
}

func (s *server) listen() error {
	return s.app.Listen(s.v.GetString("addr"))
}

func Start(db *gorm.DB, store store.Store) error {
	srv := &server{
		app:   fiber.New(),
		v:     viper.GetViper(),
		db:    db,
		store: store,
	}

	srv.setupRoutes()
	return srv.listen()
}
