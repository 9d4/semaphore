package store

import (
	"github.com/9d4/semaphore/oauth"
	"github.com/9d4/semaphore/user"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type Store struct {
	User     user.Store
	OauthApp oauth.AppStore
}

func NewStore(db *gorm.DB, rdb *redis.Client) Store {
	s := Store{
		User:     user.NewStore(db),
		OauthApp: oauth.NewStore(db, rdb),
	}

	return s
}
