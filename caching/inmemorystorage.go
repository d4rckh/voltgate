package caching

import (
	"net/http"
	"sync"
	"time"
)

// CachedResponse represents a stored HTTP response
type CachedResponse struct {
	Status    int
	Header    http.Header
	Body      []byte
	ExpiresAt time.Time
}

// InMemoryCacherStorage provides an in-memory cache with TTL
type InMemoryCacherStorage struct {
	mu    sync.RWMutex
	cache map[string]CachedResponse
}

// GetRequest retrieves a cached response if it's still valid
func (r *InMemoryCacherStorage) GetRequest(method string, cacheKey string) (int, http.Header, []byte, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fullKey := method + ":" + cacheKey
	entry, exists := r.cache[fullKey]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return 0, nil, nil, false
	}

	// Return a copy of headers to prevent modification
	headerCopy := make(http.Header)
	for k, v := range entry.Header {
		headerCopy[k] = append([]string{}, v...)
	}

	return entry.Status, headerCopy, entry.Body, true
}

// CacheRequest stores a response in memory with a TTL
func (r *InMemoryCacherStorage) CacheRequest(method string, cacheKey string, status int, header http.Header, data []byte, ttl time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullKey := method + ":" + cacheKey
	r.cache[fullKey] = CachedResponse{
		Status:    status,
		Header:    header,
		Body:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// MakeInMemoryCacherStorage initializes a new cache storage
func MakeInMemoryCacherStorage() *InMemoryCacherStorage {
	return &InMemoryCacherStorage{
		cache: make(map[string]CachedResponse),
	}
}
