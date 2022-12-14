package cmd

import (
	"context"
	"fmt"
	"github.com/9d4/semaphore/server"
	"github.com/9d4/semaphore/store"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type bootData struct {
	db     *gorm.DB
	rdb    *redis.Client
	config *server.Config
}

type (
	cobraFunc func(cmd *cobra.Command, args []string)
	bootFunc  func(cmd *cobra.Command, args []string, passData *bootData)
)

func boot(fn bootFunc) cobraFunc {
	return func(cmd *cobra.Command, args []string) {
		config := server.ParseViper(v)

		// connect db and something else here
		data := &bootData{}

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost,
			config.DBPort,
			config.DBUsername,
			config.DBPassword,
			config.DBName,
		)
		db, err := gorm.Open(postgres.Open(dsn))
		if err != nil {
			jww.FATAL.Fatal(err)
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     config.RedisAddress,
			Username: config.RedisUsername,
			Password: config.RedisPassword,
		})

		if err = rdb.Ping(context.Background()).Err(); err != nil {
			jww.FATAL.Fatal(err)
		}

		// build store
		data.db = db
		data.rdb = rdb

		// auto migrate
		jww.INFO.Print("Auto Migrating...")
		store.MigrateAll(db)
		jww.INFO.Print("done.")

		fn(cmd, args, data)
	}
}
