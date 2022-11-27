package cmd

import (
	"crypto/rand"
	"fmt"

	"github.com/9d4/semaphore/store"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	v "github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type bootData struct {
	db    *gorm.DB
	store store.Store
}

type (
	cobraFunc func(cmd *cobra.Command, args []string)
	bootFunc  func(cmd *cobra.Command, args []string, passData *bootData)
)

func boot(fn bootFunc) cobraFunc {
	return func(cmd *cobra.Command, args []string) {
		if v.GetBool("gen-key") {
			fmt.Println(generateKey())
			return
		}

		// connect db and something else here
		data := &bootData{}

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			v.GetString("dbhost"),
			v.GetInt("dbport"),
			v.GetString("dbuser"),
			v.GetString("dbpasswd"),
			v.GetString("dbname"),
		)
		db, err := gorm.Open(postgres.Open(dsn))
		if err != nil {
			jww.FATAL.Fatal(err)
		}
		data.db = db

		// build store
		data.store = store.NewStore(db)

		// auto migrate
		jww.INFO.Print("Auto Migrating...")
		store.MigrateAll(db)
		jww.INFO.Print("done.")

		// seed?
		if v.GetBool("seed") {
			jww.INFO.Print("Seeding database...")
			store.Seed(data.store)
			jww.INFO.Print("done.")
		}

		fn(cmd, args, data)
	}
}

func generateKey() string {
	buff := make([]byte, 32)
	if _, err := rand.Read(buff); err != nil {
		jww.FATAL.Fatal(err)
	}

	return fmt.Sprintf("%x", buff)
}
