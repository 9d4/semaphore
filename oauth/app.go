package oauth

import (
	"gorm.io/gorm"
	"time"
)

type App struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `json:"name"`
	ClientID     string         `json:"client_id" gorm:"uniqueIndex"`
	ClientSecret string         `json:"-"`
}
