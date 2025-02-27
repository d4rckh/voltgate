package ratelimiting

import (
	"context"
	"strconv"
	"time"
	"voltgate-proxy/storage"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiterStorage is a Redis-backed implementation of RateLimiterStorage.
type RedisRateLimiterStorage struct {
	client *redis.Client
	ctx    context.Context
}

// IncreaseRequestCount increments the request count and returns the new count and remaining TTL.
func (r *RedisRateLimiterStorage) IncreaseRequestCount(key string, window time.Duration) (int, time.Duration) {
	value, ttl := r.GetRequestCount(key, window)

	r.client.Set(r.ctx, key, value+1, redis.KeepTTL)

	return value + 1, ttl
}

// GetRequestCount retrieves the request count and remaining TTL.
func (r *RedisRateLimiterStorage) GetRequestCount(key string, window time.Duration) (int, time.Duration) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		val = "0"
		//log.Printf("[GetRequestCount] Error getting value for key %s: %v, creating the key with value 0", key, err)
		r.client.Set(r.ctx, key, 0, window)
	}

	ttl, _ := r.client.TTL(r.ctx, key).Result()

	intValue, _ := strconv.Atoi(val)

	if err != nil {
		//log.Printf("[GetRequestCount] Error parsing string value '%s' for key %s: %v", val, key, err)
		intValue = 0
	}

	if ttl > window {
		//log.Printf("Invalid TTL detected for key %s", key)
		r.client.Set(r.ctx, key, intValue, window)
	}

	return intValue, ttl
}

// ResetRequestCount resets the request count for a given key.
func (r *RedisRateLimiterStorage) ResetRequestCount(key string) {
	r.client.Del(r.ctx, key)
}

// MakeRedisRateLimiterStorage creates an instance of RedisRateLimiterStorage.
func MakeRedisRateLimiterStorage(config storage.RedisAppConfig) *RedisRateLimiterStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		Username: config.Username,
	})
	return &RedisRateLimiterStorage{
		client: client,
		ctx:    context.Background(),
	}
}
