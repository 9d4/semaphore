package middleware

import (
	"context"
	"github.com/9d4/semaphore/auth"
	"github.com/9d4/semaphore/oauth"
	"github.com/9d4/semaphore/server/types"
	"github.com/9d4/semaphore/server/util"
	"github.com/gofiber/fiber/v2"
)

func OAuthBearerAuth(key []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := util.GetBearerToken(c)
		if err != nil {
			return err
		}

		at, err := oauth.ValidateAccessToken(token, auth.DefaultJwtKeyFunc(key))
		if err != nil {
			return fiber.ErrUnauthorized
		}

		ctx := context.WithValue(c.UserContext(), types.ContextKey("access_token"), at)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
