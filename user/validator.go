package user

import (
	"github.com/go-playground/validator/v10"
)

var defaultValidate *Validator

func init() {
	if defaultValidate == nil {
		defaultValidate = &Validator{
			validate: validator.New(),
		}
	}
}

type Validator struct {
	validate *validator.Validate
}

func GetValidate() *validator.Validate { return defaultValidate.GetValidate() }

func (v *Validator) GetValidate() *validator.Validate {
	return v.validate
}

func Validate(user *User) error { return defaultValidate.Validate(user) }

func (v *Validator) Validate(user *User) error {
	return v.validate.Struct(user)
}
