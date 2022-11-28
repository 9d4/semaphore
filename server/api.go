package server

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/9d4/semaphore/store"
	"github.com/9d4/semaphore/user"
	"github.com/9d4/semaphore/util"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type apiServer struct {
	app   *fiber.App
	db    *gorm.DB
	store store.Store
	v     *viper.Viper
}

type userInfo struct {
	ID        uint
	Email     string
	FirstName string
	LastName  string
}

type accessToken struct {
	User userInfo
	jwt.RegisteredClaims
}

type refreshToken struct {
	jwt.RegisteredClaims
}

const (
	AccessTokenExpirationTime  = time.Minute * 15
	RefreshTokenExpirationTime = time.Hour * 48
)

func newApiServer(db *gorm.DB, store store.Store) *apiServer {
	srv := &apiServer{
		app:   fiber.New(),
		v:     viper.GetViper(),
		db:    db,
		store: store,
	}

	srv.setupRoutes()
	return srv
}

func (s *apiServer) setupRoutes() {
	s.app.Post("/login", s.handleLogin)
	s.app.Post("/renew", s.handleRenew)
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
		return writeError(c, ErrCredentialNotFound)
	}

	match := util.VerifyEncoded([]byte(cred.Password), []byte(usr.Password))
	if !match {
		return writeError(c, ErrCredentialNotFound)
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

	token, err := jwt.ParseWithClaims(rtRaw, &rt, jwtKeyFunc([]byte(s.v.GetString("app_key"))))
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	subjectID, err := strconv.Atoi(rt.Subject)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var usr user.User
	result := s.db.First(&usr, user.User{Model: gorm.Model{ID: uint(subjectID)}})
	if result.Error != nil {
		return fiber.ErrUnauthorized
	}

	tokenPair, err := s.generateTokenPair(usr)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.SendStatus(fiber.StatusCreated)
	return c.JSON(tokenPair)
}

func (s *apiServer) generateTokenPair(usr user.User) (map[string]interface{}, error) {
	key := []byte(s.v.GetString("app_key"))

	at, err := generateAccessToken(usr, key, AccessTokenExpirationTime)
	if err != nil {
		jww.TRACE.Println("apiServer:error:generateAccessToken", err)
		return nil, err
	}

	rt, err := generateRefreshToken(usr, key, RefreshTokenExpirationTime)
	if err != nil {
		jww.TRACE.Println("apiServer:error:generateAccessToken", err)
		return nil, err
	}

	return map[string]interface{}{
		"access_token":  at,
		"refresh_token": rt,
	}, nil
}

func generateAccessToken(usr user.User, key []byte, expiresIn time.Duration) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "semaphore",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
		User: userInfo{
			ID:        usr.ID,
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
		},
	})

	return at.SignedString(key)
}

func generateRefreshToken(usr user.User, key []byte, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{}
	claims.Subject = fmt.Sprint(usr.ID)
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiresIn))

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return rt.SignedString(key)
}

func jwtKeyFunc(key []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		return key, nil
	}
}