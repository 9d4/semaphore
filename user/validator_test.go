package user

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		wantErr        bool
		wantFieldError string
	}{
		{
			name: "Full",
			user: &User{
				Email:     "test@gg.gg",
				FirstName: "abc",
				LastName:  "xyz",
				Password:  "secret",
			},
		},
		{
			name: "Firstname 2 letters",
			user: &User{
				Email:     "test@gg.gg",
				FirstName: "as",
				LastName:  "xyz",
				Password:  "secret",
			},
			wantErr:        true,
			wantFieldError: "FirstName",
		},
		{
			name: "Password 4 letters",
			user: &User{
				Email:     "test@gg.gg",
				FirstName: "abcd",
				LastName:  "xyz",
				Password:  "secr",
			},
			wantErr:        true,
			wantFieldError: "Password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.user)
			if err != nil && tt.wantErr {
				validationErrors := err.(validator.ValidationErrors)
				validError := false

				for _, validationError := range validationErrors {
					if validationError.Field() == tt.wantFieldError {
						validError = true
					}
				}

				if !validError {
					t.Errorf("Wanted error in field %s. Got: %s", tt.wantFieldError, validationErrors.Error())
				}
			}
		})
	}
}
