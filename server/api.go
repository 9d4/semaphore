package server

import (
	"context"
	"errors"
	"github.com/9d4/semaphore/auth"
	errs "github.com/9d4/semaphore/errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/9d4/semaphore/user"
	"github.com/9d4/semaphore/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type apiServer struct {
	*Config
	app *fiber.App
	db  *gorm.DB
	v   *viper.Viper
}

type userInfo struct {
	ID        uint   `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
}

type accessToken struct {
	User userInfo
	jwt.RegisteredClaims
}

type refreshToken struct {
	jwt.RegisteredClaims
}

type contextKey string

func newApiServer(db *gorm.DB, opts ...Option) *apiServer {
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

	srv := &apiServer{
		Config: config,
		app:    fiber.New(),
		v:      viper.GetViper(),
		db:     db,
	}

	srv.setupRoutes()
	return srv
}

func (s *apiServer) setupRoutes() {
	s.app.Post("/login", s.handleLogin)
	s.app.Post("/renew", s.handleRenew)
	users := s.app.Group("users/")
	users.Get(":userid/profile", s.withAuth, s.handleUsersProfile)
}

func (s *apiServer) handleLogin(c *fiber.Ctx) error {
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

	match := util.VerifyEncoded([]byte(cred.Password), []byte(usr.Password))
	if !match {
		return errs.WriteErrorJSON(c, errs.ErrCredentialNotFound)
	}

	// if url contains query "check=1" then don't generate token
	if c.Query("check") == "1" {
		return c.SendStatus(200)
	}

	tokenPair, err := s.generateTokenPair(usr)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(tokenPair)
}

// This endpoint is not truly rest. It takes cookie containing
// refresh token set by server.
func (s *apiServer) handleRenew(c *fiber.Ctx) error {
	rtRaw := c.Cookies("rt")
	if rtRaw == "" {
		return fiber.ErrUnauthorized
	}

	var rt refreshToken

	token, err := jwt.ParseWithClaims(rtRaw, &rt, auth.DefaultJwtKeyFunc(s.KeyBytes))
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	subjectID, err := strconv.Atoi(rt.Subject)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var usr user.User
	result := s.db.First(&usr, user.User{ID: uint(subjectID)})
	if result.Error != nil {
		return fiber.ErrUnauthorized
	}

	tokenPair, err := s.generateTokenPair(usr)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.SendStatus(fiber.StatusCreated)
	c.Cookie(&fiber.Cookie{
		Name:     "rt",
		Value:    tokenPair["refresh_token"],
		Domain:   s.v.GetString("cookie_domain"),
		Expires:  time.Now().Add(auth.RefreshTokenExpiration),
		HTTPOnly: true,
	})

	return c.JSON(tokenPair)
}

func (s *apiServer) handleUsersProfile(c *fiber.Ctx) error {
	paramUserID := c.Params("userid")
	userid, err := strconv.Atoi(paramUserID)
	if err != nil {
		return fiber.ErrBadRequest
	}

	at := new(accessToken)
	at, ok := c.UserContext().Value(contextKey("access_token")).(*accessToken)
	if !ok {
		return fiber.ErrInternalServerError
	}

	if at.User.ID != uint(userid) {
		return fiber.ErrForbidden
	}

	var usr user.User
	result := s.db.First(&usr, user.User{ID: at.User.ID})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errs.WriteErrorJSON(c, errs.ErrCredentialNotFound)
	}

	return c.JSON(usr)
}

func (s *apiServer) withAuth(c *fiber.Ctx) error {
	authorizationPrefix := "Bearer "
	authHeader := c.GetReqHeaders()[fiber.HeaderAuthorization]

	if authHeader == "" {
		return fiber.ErrUnauthorized
	}

	token := strings.TrimPrefix(authHeader, authorizationPrefix)
	at, err := auth.ValidateAccessToken(token, auth.DefaultJwtKeyFunc(s.KeyBytes))
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx := context.WithValue(context.Background(), contextKey("access_token"), at)
	c.SetUserContext(ctx)
	return c.Next()
}

func (s *apiServer) generateTokenPair(usr user.User) (map[string]string, error) {
	at, rt, err := auth.GenerateTokenPair(usr, s.KeyBytes)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  at,
		"refresh_token": rt,
	}, nil
}

// Deprecated: use auth.ValidateAccessToken
func validateAccessToken(token string, key []byte) (*accessToken, error) {
	claims := accessToken{}

	tk, err := jwt.ParseWithClaims(token, &claims, auth.DefaultJwtKeyFunc(key))
	if err != nil || !tk.Valid {
		return nil, err
	}

	return &claims, nil
}
