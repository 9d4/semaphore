package oauth

import (
	"gorm.io/gorm"
)

type App struct {
	gorm.Model
	Name     string `json:"name"`
	ClientID string `json:"client_id" gorm:"unique,index"`
}
