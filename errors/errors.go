package errors

import (
	"github.com/gofiber/fiber/v2"
)

// Error represents error during request handling
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

func WriteErrorJSON(c *fiber.Ctx, err *Error) error {
	er := c.SendStatus(err.Code)
	if er != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(map[string]interface{}{
		"error":   err.ErrorID,
		"message": err.Message,
	})
}

// Errors
var (
	ErrCredentialNotFound = NewError(fiber.StatusUnauthorized, "auth_failed", "Credential not found")

	ErrOauthClientNotFound = NewError(fiber.StatusNotFound, "oauth_client_not_found", "Client not found")
)
