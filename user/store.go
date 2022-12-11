package user

import (
	"gorm.io/gorm"
	"strings"
)

type Store interface {
	// Create inserts a new user into the database.
	// Returns an error if the user could not be created.
	Create(u *User) error

	// User gets a user from the database that matches the specified criteria.
	// Returns the user and any error that occurred.
	User(u *User) (*User, error)

	// UserByID gets a user from the database with the specified ID.
	// Returns the user and any error that occurred.
	UserByID(id uint) (*User, error)

	// UserByEmail gets a user from the database with the specified email address.
	// Returns the user and any error that occurred.
	UserByEmail(email string) (*User, error)

	// UserByName gets users from the database with the specified name.
	// Returns the users and any error that occurred.
	UserByName(name string) ([]*User, error)

	// DB gets the underlying *gorm.DB instance.
	DB() *gorm.DB

	// SetDB sets the underlying *gorm.DB instance.
	SetDB(db *gorm.DB)

	// Migrate auto-migrates the User model to database.
	Migrate() error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &store{db: db}
}

func (s *store) User(u *User) (*User, error) {
	// Get user from db with custom query
	var usr User
	tx := s.db.Where(u).First(&usr)

	if tx.Error != nil {
		return nil, resolveError(tx.Error)
	}

	return &usr, nil
}

func (s *store) UserByID(id uint) (*User, error) {
	var usr User
	tx := s.db.Where("id = ?", id).First(&usr)

	if tx.Error != nil {
		return nil, resolveError(tx.Error)
	}

	return &usr, nil
}

func (s *store) UserByEmail(email string) (*User, error) {
	var usr User
	tx := s.db.Where("email = ?", email).First(&usr)

	if tx.Error != nil {
		switch err := tx.Error; {
		case err == gorm.ErrRecordNotFound:
			return nil, ErrUserNotFound
		}

		return nil, tx.Error
	}

	return &usr, nil
}

func (s *store) UserByName(name string) ([]*User, error) {
	name = strings.ToLower(name)
	var users []*User

	tx := s.db.
		Where("lower(first_name) LIKE ?", "%"+name+"%").
		Or("lower(last_name) LIKE ?", "%"+name+"%").
		Find(&users)

	if tx.Error != nil {
		return nil, resolveError(tx.Error)
	}

	return users, nil
}

func (s *store) DB() *gorm.DB {
	return s.db
}

func (s *store) SetDB(db *gorm.DB) {
	s.db = db
}

func (s *store) Migrate() error {
	err := s.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}

func (s *store) Create(u *User) error {
	tx := s.db.Create(u)

	return tx.Error
}
