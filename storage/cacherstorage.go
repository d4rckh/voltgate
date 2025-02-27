package storage

import (
	"net/http"
	"time"
)

type CacherStorage interface {
	// GetRequest gets the request by method and path, returns status code, response body and exists boolean
	GetRequest(method string, cacheKey string) (int, http.Header, []byte, bool)

	// CacheRequest caches the request for future retrieval with a time to live
	CacheRequest(method string, cacheKey string, status int, header http.Header, data []byte, ttl time.Duration)
}
