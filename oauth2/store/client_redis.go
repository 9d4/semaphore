package store

import (
	"bytes"
	"context"
	"encoding/json"
	oauth2 "github.com/9d4/semaphore/oauth2"
	"github.com/9d4/semaphore/oauth2/models"
	"github.com/go-redis/redis/v9"
)

const ClientStoreCachePrefix = "oauth:client:"

// NewClientStoreRedis create client store
func NewClientStoreRedis(rdb *redis.Client) *ClientStoreRedis {
	return &ClientStoreRedis{rdb: rdb}
}

// ClientStoreRedis client information store
type ClientStoreRedis struct {
	rdb *redis.Client
}

// GetByID according to the ID for the client information
func (cs *ClientStoreRedis) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	key := ClientStoreCachePrefix + id
	res := cs.rdb.Get(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}

	cli := models.Client{}
	err := json.NewDecoder(bytes.NewBufferString(res.Val())).Decode(&cli)
	if err != nil {
		return nil, err
	}

	return &cli, nil
}

// Set set client information
func (cs *ClientStoreRedis) Set(id string, cli oauth2.ClientInfo) (err error) {
	key := ClientStoreCachePrefix + id

	jsonBuf := &bytes.Buffer{}
	err = json.NewEncoder(jsonBuf).Encode(cli)
	if err != nil {
		return
	}

	status := cs.rdb.Set(context.Background(), key, jsonBuf.String(), redis.KeepTTL)
	if status.Err() != nil {
		err = status.Err()
		return
	}

	return
}
