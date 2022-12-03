package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `json:"email" gorm:"index:email_index,unique"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"-"`
}
