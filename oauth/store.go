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
	GenerateAuthorizationCode(scopes []Scope, clientID string, subject string) (*AuthorizationCode, error)
	GetAuthorizationCode(code string) (*AuthorizationCode, error)
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
func (s appStore) GenerateAuthorizationCode(scopes []Scope, clientID string, subject string) (*AuthorizationCode, error) {
	var (
		generated       = generateAuthCode(scopes, clientID, subject)
		generatedJson   = bytes.NewBuffer([]byte{})
		cacheKey        = fmt.Sprint(CachePrefixAuthorizationCode, generated.Code)
		cacheExpiration = time.Minute * 5
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

// GetAuthorizationCode takes marshaled AuthorizationCode from redis cache
// therefore if key does not found it returns error with type redis.Nil
func (s appStore) GetAuthorizationCode(code string) (*AuthorizationCode, error) {
	var (
		authorizationCode = AuthorizationCode{}
		key               = fmt.Sprint(CachePrefixAuthorizationCode, code)
		bgCtx, _          = context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	)

	res := s.rdb.Get(bgCtx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}

	err := json.NewDecoder(bytes.NewBufferString(res.Val())).Decode(&authorizationCode)
	if err != nil {
		return nil, err
	}

	return &authorizationCode, nil
}

func generateAuthCode(scopes []Scope, clientID string, subject string) *AuthorizationCode {
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
	authorizationCode.Subject = subject

	return authorizationCode
}
