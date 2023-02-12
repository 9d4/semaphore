package store

import (
	"github.com/9d4/semaphore/user"
	"github.com/9d4/semaphore/util"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, rdb *redis.Client) error {
	userStore := user.NewStore(db)
	userStore.Create(&user.User{
		Email:     "admin@example.com",
		FirstName: "Admin",
		Password:  hashPasswd("adm1n"),
	})
	
	return nil
}

func hashPasswd(pass string) string {
	var passHashed string
	for {
		p, err := util.HashString([]byte(pass))
		if err == nil {
			passHashed = p
			break
		}
	}

	return passHashed
}
