package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

const (
	DBLogName = "semaphore.db.log"
)

// ConnectDB returns *gorm.DB
func ConnectDB(config *Config) (*gorm.DB, error) {
	dbLogFile, err := os.OpenFile("semaphore.db.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	dbLogger := logger.New(
		log.New(dbLogFile, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ConnectRDB connects to Redis Database
func ConnectRDB(config *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
