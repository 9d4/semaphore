package util

import (
	"github.com/9d4/semaphore/server/types"
	"github.com/gofiber/fiber/v2"
	"os"
)

func UseContext[T interface{}](c *fiber.Ctx, key types.ContextKey) (T, error) {
	var thing T
	thing, ok := c.UserContext().Value(key).(T)
	if !ok {
		return thing, os.ErrNotExist
	}

	return thing, nil
}
