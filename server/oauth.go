package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/9d4/semaphore/oauth"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	jww "github.com/spf13/jwalterweatherman"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type oauthServer struct {
	app *fiber.App
	db  *gorm.DB
	rdb *redis.Client
}

type oauthScope string

// Scopes
const (
	ScopeUserinfoRead oauthScope = "ur"
)

var oauthScopes = []oauthScope{
	ScopeUserinfoRead,
}

const (
	// RedisOauthClientPrefix should be followed with clientID.
	RedisOauthClientPrefix = "oauth:client:"
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
	cacheKey := fmt.Sprint(RedisOauthClientPrefix, queryClientID)
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
	scopes := strings.Split(queryScope, "\x20")

	// if a scopes does neither not found nor malformed nor anything
	// we don't know, just remove from slice
	var fixedScopes []oauthScope
	for _, s := range scopes {
		for _, definedScope := range oauthScopes {
			// 		queried scopes available	 AND  no duplicate
			if oauthScope(s) == definedScope && !slices.Contains(fixedScopes, definedScope) {
				fixedScopes = append(fixedScopes, definedScope)
			}
		}
	}

	// Tells frontend from which backend part is the redirection.
	queryRedirectedFrom := "from=oauth_authorize"

	authorizationViewRoute := fmt.Sprint("/o", string(c.Request().URI().Path()), "?", queryRedirectedFrom, "&", string(c.Request().URI().QueryString()))
	return c.Redirect(authorizationViewRoute, fiber.StatusSeeOther)
}

type authorizationCode struct {
	Code     string       `json:"code"`
	Scopes   []oauthScope `json:"scopes"`
	ClientID string       `json:"client_id"`
}

func generateAuthorizationCode(scopes []oauthScope, clientID string) *authorizationCode {
	ac := &authorizationCode{}
	ac.Code = uuid.NewString()
	ac.Scopes = scopes
	ac.ClientID = clientID
	return ac
}
