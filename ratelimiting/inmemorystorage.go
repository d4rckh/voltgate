package ratelimiting

import (
	"sync"
	"time"
)

// InMemoryEntry represents a request entry in memory.
type InMemoryEntry struct {
	lastReset time.Time
	count     int
}

// InMemoryRateLimiterStorage is an in-memory implementation of RateLimiterStorage.
type InMemoryRateLimiterStorage struct {
	mu   sync.Mutex
	keys map[string]*InMemoryEntry
}

// IncreaseRequestCount increments the request count and returns the new count and remaining TTL.
func (i *InMemoryRateLimiterStorage) IncreaseRequestCount(key string, window time.Duration) (int, time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()

	kv, exists := i.keys[key]

	if !exists || time.Since(kv.lastReset) > window {
		// Reset if key doesn't exist or window expired
		kv = &InMemoryEntry{lastReset: time.Now(), count: 1}
		i.keys[key] = kv
	} else {
		kv.count++
	}

	return kv.count, window - time.Since(kv.lastReset)
}

// GetRequestCount retrieves the request count and remaining TTL.
func (i *InMemoryRateLimiterStorage) GetRequestCount(key string, window time.Duration) (int, time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()

	kv, exists := i.keys[key]
	if !exists || time.Since(kv.lastReset) > window {
		return 0, 0
	}

	return kv.count, window - time.Since(kv.lastReset)
}

// ResetRequestCount resets the request count for a given key.
func (i *InMemoryRateLimiterStorage) ResetRequestCount(key string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if kv, exists := i.keys[key]; exists {
		kv.count = 0
		kv.lastReset = time.Now()
	}
}

// MakeInMemoryRateLimiterStorage creates an instance of InMemoryRateLimiterStorage.
func MakeInMemoryRateLimiterStorage() *InMemoryRateLimiterStorage {
	return &InMemoryRateLimiterStorage{
		keys: make(map[string]*InMemoryEntry),
	}
}
