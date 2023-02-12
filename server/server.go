package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/9d4/semaphore/auth"
	errs "github.com/9d4/semaphore/errors"
	"github.com/9d4/semaphore/util"
	"github.com/go-redis/redis/v9"
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"github.com/9d4/semaphore/user"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
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
	requestLogWriter := jww.TRACE.Writer()
	if s.LogRequest {
		requestLogWriter = jww.INFO.Writer()
	}

	s.app.Use(fiberlogger.New(fiberlogger.Config{
		CustomTags: map[string]fiberlogger.LogFunc{
			"ips": func(output fiberlogger.Buffer, c *fiber.Ctx, data *fiberlogger.Data, extraParam string) (int, error) {
				return output.WriteString(strings.Join(c.IPs(), ">>"))
			},
		},
		Format:     "${time} ${pid} ${locals:requestid} [${ips}] [${ip}]:${port} ${status} - ${method} ${path}\n",
		Output:     requestLogWriter,
		TimeFormat: "2006/01/02 15:04:05",
	}))

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

	dbLogFile, err := os.OpenFile("semaphore.db.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println("Unable to create log file:", err)
	}
	dbLogger := logger.New(
		log.New(dbLogFile, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
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

	oauthSrv := newOauthServer(db, rdb, config)

	_srvErr := make(chan error, 1)
	srvErr = _srvErr

	_oauthSrvErr := make(chan error, 1)
	oauthSrvErr = _oauthSrvErr

	go func() { _srvErr <- srv.listen() }()
	go func() { _oauthSrvErr <- oauthSrv.Listen() }()

	return
}
