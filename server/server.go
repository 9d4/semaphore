package server

import (
	"errors"
	"fmt"
	"github.com/9d4/semaphore/auth"
	errs "github.com/9d4/semaphore/errors"
	"github.com/9d4/semaphore/util"
	"github.com/go-redis/redis/v9"
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/driver/postgres"
	"sort"
	"time"

	"github.com/9d4/semaphore/user"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type server struct {
	*Config
	app *fiber.App
	db  *gorm.DB
	rdb *redis.Client
	v   *viper.Viper
}

func (s *server) setupRoutes() {
	s.app.Static("/", "./views/dist/", fiber.Static{
		Compress: true,
		Browse:   false,
	})

	authRouter := s.app.Group("/auth")
	authRouter.Post("/login", s.handleLogin)

	oauthSrv := newOauthServer(s.db, s.rdb)
	s.app.Mount("/oauth", oauthSrv.app)

	apiSrv := newApiServer(s.db, s.Config)
	s.app.Mount("/api", apiSrv.app)

	// This is kinda tricky. Mounts will be executed lastly.
	// So if nothing found, fallback to index.html.
	s.app.Mount("/*", viewApp())
}

func (s *server) listen() error {
	return s.app.Listen(s.Address)
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
		return errs.WriteErrorJSON(c, errs.ErrCredentialNotFound)
	}

	if !util.VerifyEncoded([]byte(cred.Password), []byte(usr.Password)) {
		return errs.WriteErrorJSON(c, errs.ErrCredentialNotFound)
	}

	rt, err := auth.GenerateRefreshToken(usr, s.KeyBytes, auth.RefreshTokenExpiration)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "rt",
		Value:    rt,
		Domain:   s.v.GetString("cookie_domain"),
		Expires:  time.Now().Add(auth.RefreshTokenExpiration),
		HTTPOnly: true,
	})

	return c.RedirectBack(c.GetReqHeaders()[fiber.HeaderReferer], 302)
}

func Start(opts ...Option) error {
	config := &Config{}

	if len(opts) < 1 {
		config = &(*defaultConfig)
	}

	sort.Slice(opts, func(i, j int) bool {
		_, isConfig := opts[i].(*Config)
		_, isConfig2 := opts[j].(*Config)
		return isConfig && !isConfig2
	})

	for _, opt := range opts {
		if opt != nil {
			if applyErr := opt.Apply(config); applyErr != nil {
				return applyErr
			}
		}
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		jww.FATAL.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
	})

	srv := &server{
		Config: config,
		app:    fiber.New(),
		v:      config.v,
		db:     db,
		rdb:    rdb,
	}

	srv.setupRoutes()
	return srv.listen()
}
