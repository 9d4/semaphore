package store

import (
	"github.com/9d4/semaphore/user"
	"gorm.io/gorm"
)

type Store struct {
	User user.Store
}

func NewStore(db *gorm.DB) Store {
	s := Store{
		User: user.NewStore(db),
	}

	return s
}
