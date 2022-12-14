package oauth

import (
	"github.com/9d4/semaphore/user"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CachePrefix string

const (
	CachePrefixAuthorizationCode CachePrefix = "oauth:auth-code:"
	// CachePrefixOauthClient should be followed with clientID.
	CachePrefixOauthClient = "oauth:client:"
)

type Scope string

const (
	ScopeUserinfoRead Scope = "ur"
)

var Scopes = []Scope{
	ScopeUserinfoRead,
}

const AccessTokenIssuer = "semaphore-oauth"

type AuthorizationCode struct {
	Code     string  `json:"code"`
	Scopes   []Scope `json:"scopes"`
	ClientID string  `json:"client_id"`
	Subject  string  `json:"subject"`
}

type AccessToken struct {
	jwt.RegisteredClaims
	User UserInfo `json:"user"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func GenerateAccessToken(usr user.User, key []byte, expiresIn time.Duration) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    AccessTokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
		User: UserInfo{
			ID:        usr.ID,
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
		},
	})

	return at.SignedString(key)
}

func ValidateAccessToken(token string, keyFunc jwt.Keyfunc) (*AccessToken, error) {
	claims := AccessToken{}

	tk, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil || !tk.Valid {
		return nil, err
	}

	// check if claims is created from oauth, not from  normal api
	if claims.Issuer != AccessTokenIssuer {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims, nil
}
