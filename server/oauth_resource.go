package server

import (
	"errors"
	"github.com/9d4/semaphore/oauth2/generates"
	"github.com/9d4/semaphore/server/middleware"
	"github.com/9d4/semaphore/server/types"
	"github.com/9d4/semaphore/user"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"os"
	"sort"
)

type oAuthResourceServer struct {
	*Config
	*fiber.App
	db        *gorm.DB
	userStore user.Store
}

func newOAuthResourceServer(db *gorm.DB, opts ...Option) *oAuthResourceServer {
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
			opt.Apply(config)
		}
	}

	srv := &oAuthResourceServer{
		Config: config,
		App:    fiber.New(),
		db:     db,
	}
	srv.userStore = user.NewStore(db)

	srv.setupRoutes()
	return srv
}

func (s *oAuthResourceServer) setupRoutes() {
	bearerAuth := middleware.OAuthBearerAuth(s.KeyBytes)

	router := s.Group("/", bearerAuth)
	router.Get("/userinfo", s.handleUserInfo)
}

func (s *oAuthResourceServer) handleUserInfo(c *fiber.Ctx) error {
	at, err := useContext[generates.JWTAccessClaims](c, "access_token")
	if err != nil {
		return fiber.ErrInternalServerError
	}

	usr, err := s.userStore.UserByID(cast.ToUint(at.Subject))
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	userinfo := struct {
		ID        uint           `json:"id"`
		Email     string         `json:"email"`
		FirstName string         `json:"firstname"`
		LastName  string         `json:"lastname"`
	}{
		ID: usr.ID,
		Email: usr.Email,
		FirstName : usr.FirstName,
		LastName : usr.LastName,
	}

	return c.JSON(userinfo)
}

func useContext[T interface{}](c *fiber.Ctx, key types.ContextKey) (T, error) {
	var thing T
	thing, ok := c.UserContext().Value(key).(T)
	if !ok {
		return thing, os.ErrNotExist
	}

	return thing, nil
}
