package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey"`
	Email     string         `json:"email" gorm:"index:email_index,unique"`
	FirstName string         `json:"firstname"`
	LastName  string         `json:"lastname"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
