package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type AppStore interface {
	Create(a *App) error
	GenerateAuthorizationCode(scopes []Scope, clientID string) (*AuthorizationCode, error)
}

type appStore struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewStore(db *gorm.DB, rdb *redis.Client) AppStore {
	return &appStore{db: db, rdb: rdb}
}

func (s appStore) Create(a *App) error {
	tx := s.db.Create(a)

	return tx.Error
}

// GenerateAuthorizationCode generates authorization code and save to cache
func (s appStore) GenerateAuthorizationCode(scopes []Scope, clientID string) (*AuthorizationCode, error) {
	var (
		generated       = generateAuthCode(scopes, clientID)
		generatedJson   = bytes.NewBuffer([]byte{})
		cacheKey        = fmt.Sprint(CachePrefixAuthorizationCode, generated.Code)
		cacheExpiration = time.Duration(time.Minute * 5)
	)

	err := json.NewEncoder(generatedJson).Encode(generated)
	if err != nil {
		return nil, err
	}

	status := s.rdb.Set(context.Background(), cacheKey, generatedJson.String(), cacheExpiration)
	if status.Err() != nil {
		jww.TRACE.Println("error saving authorization code to cache:", status.Err())
		return nil, status.Err()
	}
	jww.TRACE.Println("saving authorization code to cache:", cacheKey)

	return generated, nil
}

func generateAuthCode(scopes []Scope, clientID string) *AuthorizationCode {
	var (
		randomRanges         = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		randomByteRangeStart = 0
		randomByteRangeEnd   = len(randomRanges) - 1
		randomLength         = 11

		output            []rune
		authorizationCode = &AuthorizationCode{}
	)

	for i := 0; i < randomLength; i++ {
		rand.Seed(time.Now().UnixNano())
		output = append(output, randomRanges[rand.Intn(randomByteRangeEnd-randomByteRangeStart)+randomByteRangeStart])
	}

	authorizationCode.Code = string(output)
	authorizationCode.Scopes = scopes
	authorizationCode.ClientID = clientID

	return authorizationCode
}
