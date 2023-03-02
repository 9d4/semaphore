package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey"`
	UUID      string         `json:"uuid" gorm:"index:uuid_index,unique"`
	Email     string         `json:"email" gorm:"index:email_index,unique" validate:"required,email"`
	FirstName string         `json:"firstname" validate:"required,min=3"`
	LastName  string         `json:"lastname" validate:"required,min=3"`
	Password  string         `json:"-" validate:"required,min=5"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// UserFieldJsonMap represents user's struct field for json key
var UserFieldJsonMap = map[string]string{
	"Email":     "email",
	"FirstName": "firstname",
	"LastName":  "lastname",
	"Password":  "password",
}
