package util

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetBearerToken(c *fiber.Ctx) (string, error) {
	authorizationPrefix := "Bearer "
	authHeader := c.GetReqHeaders()[fiber.HeaderAuthorization]

	if authHeader == "" {
		return "", fiber.ErrUnauthorized
	}

	token := strings.TrimPrefix(authHeader, authorizationPrefix)
	return token, nil
}
