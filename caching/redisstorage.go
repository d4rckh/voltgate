package caching

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"voltgate-proxy/storage"

	"github.com/redis/go-redis/v9"
)

// RedisCacherStorage is a Redis-backed implementation of CacherStorage.
type RedisCacherStorage struct {
	client *redis.Client
	ctx    context.Context
}

func (r *RedisCacherStorage) GetRequest(method string, cacheKey string) (int, http.Header, []byte, bool) {
	body, err := r.client.Get(r.ctx, strings.Join([]string{method, cacheKey, "body"}, "")).Bytes()

	if err != nil {
		return 0, make(http.Header), nil, false
	}

	status, err := r.client.Get(r.ctx, strings.Join([]string{method, cacheKey, "status"}, "")).Int()

	if err != nil {
		return 0, nil, nil, false
	}

	marshaledHeader, err := r.client.Get(r.ctx, strings.Join([]string{method, cacheKey, "header"}, "")).Bytes()
	header := make(http.Header)

	err = json.Unmarshal(marshaledHeader, &header)
	if err != nil {
		return 0, nil, nil, false
	}

	return status, header, body, true
}

func (r *RedisCacherStorage) CacheRequest(method string, cacheKey string, status int, header http.Header, data []byte, ttl time.Duration) {
	r.client.Set(r.ctx, strings.Join([]string{method, cacheKey, "body"}, ""), data, ttl)
	r.client.Set(r.ctx, strings.Join([]string{method, cacheKey, "status"}, ""), status, ttl)

	marshaledHeader, err := json.Marshal(header)

	if err != nil {
		marshaledHeader = nil
		return
	}

	r.client.Set(r.ctx, strings.Join([]string{method, cacheKey, "header"}, ""), marshaledHeader, ttl)
}

// MakeRedisCacherStorage creates an instance of RedisCacherStorage.
func MakeRedisCacherStorage(config storage.RedisAppConfig) *RedisCacherStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		Username: config.Username,
	})

	return &RedisCacherStorage{
		client: client,
		ctx:    context.Background(),
	}
}
