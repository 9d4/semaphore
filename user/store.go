package user

import "gorm.io/gorm"

type Store interface {
	Create(u *User) error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &store{db: db}
}

func (s store) Create(u *User) error {
	tx := s.db.Create(u)

	return tx.Error
}
