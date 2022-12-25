package middleware

import (
	"context"
	"github.com/9d4/semaphore/auth"
	"github.com/9d4/semaphore/oauth2/generates"
	"github.com/9d4/semaphore/server/types"
	"github.com/9d4/semaphore/server/util"
	"github.com/9d4/semaphore/server/vars"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func OAuthBearerAuth(key []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := util.GetBearerToken(c)
		if err != nil {
			return err
		}

		claims := generates.JWTAccessClaims{}

		tk, err := jwt.ParseWithClaims(token, &claims, auth.DefaultJwtKeyFunc(key))
		if err != nil || !tk.Valid || tk.Header["kid"] != vars.OAuth2AccessTokenKID {
			return fiber.ErrUnauthorized
		}


		ctx := context.WithValue(c.UserContext(), types.ContextKey("access_token"), claims)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
