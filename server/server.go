package server

import (
	"errors"
	"github.com/9d4/semaphore/util"
	"github.com/go-redis/redis/v9"
	"time"

	"github.com/9d4/semaphore/store"
	"github.com/9d4/semaphore/user"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type server struct {
	app   *fiber.App
	db    *gorm.DB
	rdb   *redis.Client
	store store.Store
	v     *viper.Viper
}

func (s *server) setupRoutes() {
	s.app.Static("/", "./views/dist/", fiber.Static{
		Compress: true,
		Browse:   false,
	})

	auth := s.app.Group("/auth")
	auth.Post("/login", s.handleLogin)

	oauthSrv := newOauthServer(s.db, s.rdb)
	s.app.Mount("/oauth", oauthSrv.app)

	apiSrv := newApiServer(s.db, s.store)
	s.app.Mount("/api", apiSrv.app)

	// This is kinda tricky. Mounts will be executed lastly.
	// So if nothing found, fallback to index.html.
	staticSrv := newViewServer()
	s.app.Mount("/*", staticSrv.app)
}

func (s *server) listen() error {
	return s.app.Listen(s.v.GetString("addr"))
}

func (s *server) handleLogin(c *fiber.Ctx) error {
	type jsonCred struct {
		Email    string
		Password string
	}

	cred := new(jsonCred)
	if err := c.BodyParser(cred); err != nil {
		return err
	}

	var usr user.User
	result := s.db.First(&usr, user.User{Email: cred.Email})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return writeError(c, ErrCredentialNotFound)
	}

	if !util.VerifyEncoded([]byte(cred.Password), []byte(usr.Password)) {
		return writeError(c, ErrCredentialNotFound)
	}

	rt, err := generateRefreshToken(usr, []byte(s.v.GetString("app_key")), RefreshTokenExpirationTime)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "rt",
		Value:    rt,
		Domain:   s.v.GetString("cookie_domain"),
		Expires:  time.Now().Add(RefreshTokenExpirationTime),
		HTTPOnly: true,
	})

	return c.RedirectBack(c.GetReqHeaders()[fiber.HeaderReferer], 302)
}

func Start(db *gorm.DB, rdb *redis.Client, store store.Store) error {
	srv := &server{
		app:   fiber.New(),
		v:     viper.GetViper(),
		db:    db,
		rdb:   rdb,
		store: store,
	}

	srv.setupRoutes()
	return srv.listen()
}
