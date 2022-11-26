package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"index:email_index,unique"`
	FirstName string
	LastName  string
}
