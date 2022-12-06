package store

import (
	"github.com/9d4/semaphore/oauth"
	"github.com/9d4/semaphore/user"
	"github.com/9d4/semaphore/util"
)

func Seed(s Store) error {
	s.User.Create(&user.User{
		Email:     "admin@example.com",
		FirstName: "Admin",
		Password:  hashPasswd("adm1n"),
	})

	s.OauthApp.Create(&oauth.App{
		Name:     "TestApp",
		ClientID: "test-app-client_id",
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
