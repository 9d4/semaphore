package server

import (
	"context"
	"errors"
	"github.com/9d4/semaphore/auth"
	errs "github.com/9d4/semaphore/errors"
	"github.com/9d4/semaphore/server/middleware"
	"github.com/9d4/semaphore/server/types"
	"github.com/go-playground/validator/v10"
	jww "github.com/spf13/jwalterweatherman"
	"sort"
	"strconv"
	"time"

	serverutil "github.com/9d4/semaphore/server/util"
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
	bearerAuth := middleware.BearerAuth(s.KeyBytes)

	s.app.Post("/login", s.handleLogin)
	s.app.Post("/renew", s.handleRenew)
	users := s.app.Group("users/")
	users.Get(":userid/profile", bearerAuth, s.handleUsersProfile)
	users.Post("/", s.handleUsersStore)
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

	at := &auth.AccessToken{}
	at, ok := c.UserContext().Value(types.ContextKey("access_token")).(*auth.AccessToken)
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

func (s *apiServer) handleUsersStore(c *fiber.Ctx) error {
	body := struct {
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Password  string `json:"password"`
	}{}
	err := c.BodyParser(&body)
	if err != nil {
		jww.ERROR.Println("user:store:", err)
		return fiber.ErrInternalServerError
	}

	hashedPwd, err := util.HashString(util.StringToBytes(body.Password))
	if err != nil {
		jww.ERROR.Println("unable to hash password:", err)
		return fiber.ErrInternalServerError
	}

	usr := &user.User{
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Password:  hashedPwd,
	}

	err = user.Validate(usr)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return replyValidationErrors(c, validationErrors)
	}

	err = user.NewStore(s.db).Create(usr)
	if err != nil {
		jww.ERROR.Println("unable to store user on register:", err)

		return fiber.ErrInternalServerError
	}

	c.SendStatus(fiber.StatusCreated)
	return c.JSON(usr)
}

func (s *apiServer) bearerAuth(c *fiber.Ctx) error {
	token, err := serverutil.GetBearerToken(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	at, err := auth.ValidateAccessToken(token, auth.DefaultJwtKeyFunc(s.KeyBytes))
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx := context.WithValue(context.Background(), types.ContextKey("access_token"), at)
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

func replyValidationErrors(c *fiber.Ctx, validationErrors validator.ValidationErrors) error {
	type e struct {
		// represents what field has error
		Field string `json:"field"`
		Tag   string `json:"tag"`
		Param string `json:"param"`
	}
	errrs := struct {
		Errors []e `json:"errors"`
	}{}

	for _, ve := range validationErrors {
		errrs.Errors = append(errrs.Errors, e{user.UserFieldJsonMap[ve.Field()], ve.Tag(), ve.Param()})
	}

	c.SendStatus(fiber.StatusBadRequest)
	return c.JSON(errrs)
}
