package storage

import "time"

type RateLimiterStorage interface {
	// IncreaseRequestCount Returns: new request count and TTL
	IncreaseRequestCount(key string, window time.Duration) (int, time.Duration)

	// GetRequestCount Returns: request count and TTL
	GetRequestCount(key string, window time.Duration) (int, time.Duration)

	// ResetRequestCount resets the count
	ResetRequestCount(key string)
}
