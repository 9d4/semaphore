package server

import (
	"errors"
	"fmt"
	"github.com/9d4/semaphore/database"
	"github.com/9d4/semaphore/store"
	"sort"
	"strings"
	"time"

	"github.com/9d4/semaphore/auth"
	errs "github.com/9d4/semaphore/errors"
	"github.com/9d4/semaphore/user"
	"github.com/9d4/semaphore/util"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	jww "github.com/spf13/jwalterweatherman"
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
	s.app.Use(s.handleLogger())

	s.app.Static("/", "./views/dist/", fiber.Static{
		Compress: true,
		Browse:   false,
	})

	authRouter := s.app.Group("/auth")
	authRouter.Post("/login", s.handleLogin)

	oauthResourceServer := newOAuthResourceServer(s.db, s.Config)
	s.app.Mount("/api/oauth2", oauthResourceServer.App)

	apiSrv := newApiServer(s.db, s.Config)
	s.app.Mount("/api", apiSrv.app)

	// This is kinda tricky. Mounts will be executed lastly.
	// So if nothing found, fallback to index.html.
	s.app.Mount("/*", viewApp())
}

func (s *server) listen() error {
	return s.app.Listen(s.Address)
}

func (s *server) handleLogger() fiber.Handler {
	requestLogWriter := jww.TRACE.Writer()
	if s.LogRequest {
		requestLogWriter = jww.INFO.Writer()
	}

	return fiberlogger.New(fiberlogger.Config{
		CustomTags: map[string]fiberlogger.LogFunc{
			"ips": func(output fiberlogger.Buffer, c *fiber.Ctx, data *fiberlogger.Data, extraParam string) (int, error) {
				return output.WriteString(strings.Join(c.IPs(), ">>"))
			},
		},
		Format:     "${time} ${pid} ${locals:requestid} [${ips}] [${ip}]:${port} ${status} - ${method} ${path}\n",
		Output:     requestLogWriter,
		TimeFormat: "2006/01/02 15:04:05",
	})
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

func Start(opts ...Option) (srvErr <-chan error, oauthSrvErr <-chan error) {
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
				jww.FATAL.Fatal(applyErr)
			}
		}
	}

	db, err := database.ConnectDB(&database.Config{
		Host:     config.DBHost,
		Port:     config.DBPort,
		Database: config.DBName,
		Username: config.DBUsername,
		Password: config.DBPassword,
	})
	if err != nil {
		jww.FATAL.Fatalf("unable to connect to database: %s", err)
	}

	rdb, err := database.ConnectRDB(&database.RedisConfig{
		Address:  config.RedisAddress,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
		DB:       0,
	})
	if err != nil {
		jww.FATAL.Fatalf("unable to connect to redis database: %s", err)
	}

	// auto migrate
	fmt.Print("Auto Migrating...")
	store.MigrateAll(db)
	fmt.Println("\rAuto Migrating...done.")

	srv := &server{
		Config: config,
		app:    fiber.New(),
		v:      config.v,
		db:     db,
		rdb:    rdb,
	}
	srv.setupRoutes()

	oauthSrv := newOauthServer(db, rdb, config)

	_srvErr := make(chan error, 1)
	srvErr = _srvErr

	_oauthSrvErr := make(chan error, 1)
	oauthSrvErr = _oauthSrvErr

	go func() { _srvErr <- srv.listen() }()
	go func() { _oauthSrvErr <- oauthSrv.Listen() }()

	return
}
