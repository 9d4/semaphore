package server

import (
	"github.com/gofiber/fiber/v2"
)

// Represents error during request handling
type Error struct {
	Code    int    `json:"code"`
	ErrorID string `json:"error"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, errorID string, message string) *Error {
	err := &Error{
		Code:    code,
		ErrorID: errorID,
		Message: message,
	}
	return err
}

func writeError(c *fiber.Ctx, err *Error) error {
	c.SendStatus(err.Code)
	return c.JSON(map[string]interface{}{
		"error":   err.ErrorID,
		"message": err.Message,
	})
}

// Errors
var (
	ErrCredentialNotFound = NewError(fiber.StatusUnauthorized, "auth_failed", "Credential not found")
)
