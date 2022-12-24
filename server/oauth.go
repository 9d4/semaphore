package server

import (
	"errors"
	"fmt"
	"github.com/9d4/semaphore/auth"
	"github.com/9d4/semaphore/oauth2"
	"github.com/9d4/semaphore/oauth2/generates"
	"github.com/9d4/semaphore/oauth2/manage"
	"github.com/9d4/semaphore/oauth2/models"
	o2server "github.com/9d4/semaphore/oauth2/server"
	"github.com/9d4/semaphore/oauth2/store"
	oredis "github.com/9d4/semaphore/oauth2/store/redis"
	"github.com/9d4/semaphore/user"
	redis8 "github.com/go-redis/redis/v8"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

var ErrSuspended = errors.New("waiting user authorization")

type oauthServer struct {
	*Config
	app *fiber.App
	db  *gorm.DB
	rdb *redis.Client

	manager *manage.Manager
	server  *o2server.Server
	mux     *http.ServeMux
}

func newOauthServer(db *gorm.DB, rdb *redis.Client, config *Config) *oauthServer {
	os := &oauthServer{
		Config: config,
		app:    fiber.New(),
		db:     db,
		rdb:    rdb,
	}

	os.manager = manage.NewDefaultManager()
	os.manager.MapAccessGenerate(generates.NewJWTAccessGenerate("semaphore-oauth2", config.KeyBytes, jwt.SigningMethodHS512))

	// storages
	clientStore := store.NewClientStoreRedis(rdb)
	_ = clientStore.Set("mymoodle", &models.Client{
		ID:     "mymoodle",
		Secret: "mymoodle-secret",
		Domain: "moodle.test",
	})
	os.manager.MapClientStorage(clientStore)
	os.manager.MapTokenStorage(oredis.NewRedisStore(&redis8.Options{
		Addr: config.RedisAddress,
		DB:   2,
	}))

	srv := o2server.NewServer(&o2server.Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{oauth2.Code, oauth2.Token},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.Refreshing,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengePlain,
			oauth2.CodeChallengeS256,
		},
	}, os.manager)

	srv.SetClientInfoHandler(o2server.ClientBasicHandler)
	srv.SetUserAuthorizationHandler(os.handleUserAuthorization)
	srv.SetAuthorizeScopeHandler(os.handleAuthorizeScope)

	os.mux = http.NewServeMux()
	os.mux.HandleFunc("/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			if err != ErrSuspended {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}
	})
	os.mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			jww.ERROR.Println(err)
		}
	})

	os.server = srv
	return os
}

func (s *oauthServer) Listen() error {
	address := strings.Split(s.Config.Address, ":")
	newPort, err := strconv.Atoi(address[1])
	if err != nil {
		jww.FATAL.Fatal(err)
	}
	newAddr := fmt.Sprint(address[0], ":", newPort+1)

	jww.INFO.Println("OAuth Server listening on", newAddr)
	return http.ListenAndServe(newAddr, s.mux)
}

// check if user authenticated or not and consent screen
func (s *oauthServer) handleUserAuthorization(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	rtCookie, err := r.Cookie("rt")
	if err != nil {
		s.redirectConsent(w, r, "oauth_authorize")
		return "", fiber.ErrUnauthorized
	}

	var rt refreshToken

	token, err := jwt.ParseWithClaims(rtCookie.Value, &rt, auth.DefaultJwtKeyFunc(s.KeyBytes))
	if err != nil || !token.Valid {
		s.redirectConsent(w, r, "oauth_authorize")
		return "", ErrSuspended
	}

	subjectID, err := strconv.Atoi(rt.Subject)
	if err != nil {
		return "", fiber.ErrInternalServerError
	}

	var usr user.User
	result := s.db.First(&usr, user.User{ID: uint(subjectID)})
	if result.Error != nil {
		return "", fiber.ErrUnauthorized
	}

	queryConsent := r.FormValue("consent")
	if queryConsent == "1" {
		return strconv.Itoa(int(usr.ID)), nil
	}

	w.Header().Set("Location", "/o/oauth/authorize?"+r.URL.RawQuery)
	w.WriteHeader(http.StatusFound)
	return "", ErrSuspended
}

func (s *oauthServer) redirectConsent(w http.ResponseWriter, r *http.Request, from string) {
	w.Header().Set("Location", "/o/oauth/authorize?"+r.URL.RawQuery+"&from="+from)
	w.WriteHeader(http.StatusFound)
}

func (s *oauthServer) handleAuthorizeScope(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	reqScopesRaw := r.FormValue("scope")
	reqScopes := strings.Split(strings.TrimSpace(reqScopesRaw), " ")
	for _, s := range reqScopes {
		if OAuth2Scopes[s] != "" {
			scope = scope + s + " "
		}
	}

	scope = strings.TrimSpace(scope)
	return
}

type OAuth2Scope string

const (
	ScopeEmail OAuth2Scope = "email"
)

var OAuth2Scopes map[string]OAuth2Scope = getOAuth2Scopes()

func getOAuth2Scopes() map[string]OAuth2Scope {
	m := make(map[string]OAuth2Scope)
	m[string(ScopeEmail)] = ScopeEmail
	return m
}
