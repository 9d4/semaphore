package user

import (
	"errors"
	"gorm.io/gorm"
)

// Error represents errors for user packages whilst maintaining
// base error. Error can either be checked using == or by calling
// errors.Is() to check the base of Error.
//
// Example of User not found in database.
//
//	user, err := s.UserByEmail(tt.args.email)
//	if err != nil {
//		if err == ErrUserNotFound {
//			// omitted
//		}
//		//or
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			// omitted
//		}
//	}
type Error struct {
	base    error
	message string
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Is(target error) bool {
	return target == e.base
}

// New creates Error with base from other error, like from gorm.
func New(base error, msg string) *Error {
	return &Error{
		base:    base,
		message: msg,
	}
}

var (
	ErrUserNotFound = New(gorm.ErrRecordNotFound, "user not found")
)

func resolveError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrUserNotFound
	default:
		return err
	}
}
