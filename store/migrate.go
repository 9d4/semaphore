package store

import (
	"github.com/9d4/semaphore/oauth"
	"github.com/9d4/semaphore/user"
	"gorm.io/gorm"
)

func MigrateAll(db *gorm.DB) {
	toBeMigrated := []interface{}{
		&user.User{},
		&oauth.App{},
	}

	db.AutoMigrate(toBeMigrated...)
}
