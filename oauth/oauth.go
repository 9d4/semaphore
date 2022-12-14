package oauth

import (
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
}

func GenerateAccessToken(key []byte, expiresIn time.Duration, subject string, clientID ...string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    AccessTokenIssuer,
			Audience:  clientID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
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
