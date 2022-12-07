package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/9d4/semaphore/oauth"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	jww "github.com/spf13/jwalterweatherman"
	v "github.com/spf13/viper"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type oauthServer struct {
	app *fiber.App
	db  *gorm.DB
	rdb *redis.Client
}

const (
	OauthAccessTokenExpirationTime = time.Duration(time.Hour * 24)
)

func newOauthServer(db *gorm.DB, rdb *redis.Client) *oauthServer {
	os := &oauthServer{
		app: fiber.New(),
		db:  db,
		rdb: rdb,
	}
	os.setupRoutes()

	return os
}

func (s *oauthServer) setupRoutes() {
	s.app.Get("/authorize", s.handleAuthorize)
	s.app.Post("/authorize", s.withAuth, s.handleAuthorize, s.handleAuthorizePost)
	s.app.Post("/token", s.handleExchangeToken)
}

func (s *oauthServer) handleAuthorize(c *fiber.Ctx) error {
	queryResponseType := c.Query("response_type")
	if queryResponseType != "code" {
		return fiber.ErrNotAcceptable
	}

	queryResponseMode := c.Query("response_mode")
	if queryResponseMode != "query" {
		return fiber.ErrNotAcceptable
	}

	queryRedirectUri := c.Query("redirect_uri")
	if queryRedirectUri == "" {
		queryRedirectUri = c.GetReqHeaders()[fiber.HeaderReferer]
	}

	queryClientID := c.Query("client_id")
	if queryClientID == "" {
		return fiber.ErrNotAcceptable
	}

	// Get oauth app from cache first then check on database
	clientApp := &oauth.App{}
	cacheKey := fmt.Sprint(oauth.CachePrefixOauthClient, queryClientID)
	cache := s.rdb.Get(context.Background(), cacheKey)
	if cache.Err() != nil {
		tx := s.db.First(clientApp, oauth.App{ClientID: queryClientID})
		if tx.Error != nil {
			if errors.Is(gorm.ErrRecordNotFound, tx.Error) {
				return writeError(c, ErrOauthClientNotFound)
			}
			jww.TRACE.Println(tx.Error)
			return fiber.ErrNotAcceptable
		}

		// save to cache
		status := s.rdb.Set(context.Background(), cacheKey, clientApp.ID, redis.KeepTTL)
		jww.DEBUG.Println(status.Err())
		jww.DEBUG.Println("save client to cache", cacheKey)
	}

	queryScope := c.Query("scope")
	if queryScope == "" {
		return fiber.ErrNotAcceptable
	}

	var scopes []oauth.Scope
	for _, s := range strings.Split(queryScope, " ") {
		scopes = append(scopes, oauth.Scope(s))
	}

	// if a scopes does neither not found nor malformed nor anything
	// we don't know, just remove from slice
	var fixedScopes []oauth.Scope
	for _, s := range scopes {
		for _, definedScope := range oauth.Scopes {
			// 		queried scopes available	 AND  no duplicate
			if oauth.Scope(s) == definedScope && !slices.Contains(fixedScopes, definedScope) {
				fixedScopes = append(fixedScopes, definedScope)
			}
		}
	}

	// if the method is POST, it means the resource owner has authorized
	// the authorization. Now generate authorization code and redirect to
	// redirect_uri value from query.
	if c.Method() == fiber.MethodPost {
		scopesCtx := context.WithValue(c.UserContext(), contextKey("scopes"), fixedScopes)
		c.SetUserContext(scopesCtx)

		scopesCtx = context.WithValue(c.UserContext(), contextKey("scopess"), fixedScopes)
		c.SetUserContext(scopesCtx)

		return c.Next()
	}

	// Tells frontend from which backend part is the redirection.
	queryRedirectedFrom := "from=oauth_authorize"

	authorizationViewRoute := fmt.Sprint("/o", string(c.Request().URI().Path()), "?", queryRedirectedFrom, "&", string(c.Request().URI().QueryString()))
	return c.Redirect(authorizationViewRoute, fiber.StatusSeeOther)
}

func (s *oauthServer) handleAuthorizePost(c *fiber.Ctx) error {
	var (
		stores = oauth.NewStore(s.db, s.rdb)

		scopes []oauth.Scope
	)

	c.UserContext()
	scopes, ok := c.UserContext().Value(contextKey("scopess")).([]oauth.Scope)
	if !ok {
		return fiber.ErrInternalServerError
	}

	at, ok := c.UserContext().Value(contextKey("access_token")).(*accessToken)
	if !ok {
		return fiber.ErrInternalServerError
	}

	authorizationCode, err := stores.GenerateAuthorizationCode(scopes, c.Query("client_id"), strconv.Itoa(int(at.User.ID)))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	// Send redirect uri to frontend
	targetUri := fmt.Sprint(c.Query("redirect_uri"), "?code=", authorizationCode.Code)

	if c.SendStatus(fiber.StatusCreated) != nil {
		return err
	}

	return c.JSON(map[string]interface{}{
		"target_uri": targetUri,
	})
}

func (s *oauthServer) handleExchangeToken(c *fiber.Ctx) error {
	var body struct {
		Code string `json:"authorization_code"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fiber.ErrBadRequest
	}

	store := oauth.NewStore(s.db, s.rdb)
	authorizationCode, err := store.GetAuthorizationCode(body.Code)
	if err != nil {
		if err == redis.Nil {
			return fiber.ErrNotFound
		}
		jww.TRACE.Println("oauth:handleExchangeToken:error getting auth code detail from cache:", err)
		return fiber.ErrInternalServerError
	}

	at, err := s.generateAccessToken(authorizationCode.Subject, authorizationCode.ClientID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(map[string]interface{}{
		"access_token": at,
	})
}

func (s *oauthServer) withAuth(c *fiber.Ctx) error {
	authorizationPrefix := "Bearer "
	authHeader := c.GetReqHeaders()[fiber.HeaderAuthorization]

	if authHeader == "" {
		return fiber.ErrUnauthorized
	}

	token := strings.TrimPrefix(authHeader, authorizationPrefix)
	at, err := validateAccessToken(token, []byte(v.GetString("app_key")))
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx := context.WithValue(c.UserContext(), contextKey("access_token"), at)
	c.SetUserContext(ctx)
	return c.Next()
}

func (s *oauthServer) generateAccessToken(subject string, clientID string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "semaphore",
			Subject:   subject,
			Audience:  jwt.ClaimStrings{clientID},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(OauthAccessTokenExpirationTime)),
		},
	})

	return at.SignedString([]byte(v.GetString("app_key")))
}
