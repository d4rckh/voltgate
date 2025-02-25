package rate_limiting

import "time"

type InMemoryEntry struct {
	lastReset time.Time
	count     int
}

type InMemoryRateLimiterStorage struct {
	keys map[string]*InMemoryEntry
}

func (i InMemoryRateLimiterStorage) IncreaseRequestCount(key string) {
	kv, exists := i.keys[key]
	if !exists {
		i.keys[key] = new(InMemoryEntry)
		i.keys[key].lastReset = time.Now()
		i.keys[key].count = 1
		kv = i.keys[key]
	}
	kv.count++
}

func (i InMemoryRateLimiterStorage) GetRequestCount(key string) (int, time.Duration) {
	kv, exists := i.keys[key]
	if !exists {
		return 0, time.Duration(0)
	}
	return kv.count, time.Since(kv.lastReset)
}

func (i InMemoryRateLimiterStorage) ResetRequestCount(key string) {
	kv, exists := i.keys[key]
	if !exists {
		return
	}
	kv.count = 0
	kv.lastReset = time.Now()
}

func MakeInMemoryRateLimiterStorage() *InMemoryRateLimiterStorage {
	return &InMemoryRateLimiterStorage{
		keys: make(map[string]*InMemoryEntry),
	}
}
