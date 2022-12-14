package auth

import (
	"fmt"
	"github.com/9d4/semaphore/user"
	"github.com/9d4/semaphore/util"
	"github.com/golang-jwt/jwt/v4"
	jww "github.com/spf13/jwalterweatherman"
	v "github.com/spf13/viper"
	"time"
)

// AccessToken represents jwt claims for user.
type AccessToken struct {
	jwt.RegisteredClaims
	User UserInfo `json:"user"`
}

// RefreshToken represents jwt claims for user refresh token.
type RefreshToken struct {
	jwt.RegisteredClaims
}

// UserInfo represents user info that will be part of AccessToken.
type UserInfo struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

const (
	AccessTokenExpiration  = time.Minute * 15
	RefreshTokenExpiration = time.Hour * 48
)

const AccessTokenIssuer = "semaphore"

var DefaultJwtKeyFunc = func(key []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}
}

var key []byte

func GenerateAccessToken(usr user.User, key []byte, expiresIn time.Duration) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "semaphore",
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

func GenerateRefreshToken(usr user.User, key []byte, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{}
	claims.Subject = fmt.Sprint(usr.ID)
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiresIn))

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return rt.SignedString(key)
}

func GenerateTokenPair(usr user.User, key []byte) (accessToken string, refreshToken string, err error) {
	at, genErr := GenerateAccessToken(usr, key, AccessTokenExpiration)
	if genErr != nil {
		jww.TRACE.Println("apiServer:error:generateAccessToken", genErr)
		err = genErr
		return
	}

	rt, genErr := GenerateRefreshToken(usr, key, RefreshTokenExpiration)
	if genErr != nil {
		jww.TRACE.Println("apiServer:error:generateAccessToken", genErr)
		err = genErr
		return
	}

	accessToken, refreshToken = at, rt
	return
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

func ValidateRefreshToken(token string, keyFunc jwt.Keyfunc) (*RefreshToken, error) {
	claims := RefreshToken{}

	tk, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil || !tk.Valid {
		return nil, err
	}

	// check if claims is created from normal auth, not from oauth
	if claims.Issuer != AccessTokenIssuer {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims, nil
}

func getKey() []byte {
	k := ""
	if key == nil {
		k = v.GetString("app_key")
		key = util.StringToBytes(k)
	}

	return key
}
