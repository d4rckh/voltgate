package proxy

import (
	"net/http"
	"net/url"
	"sync"
	"voltgate-proxy/storage"
)

type Server struct {
	Mu                 sync.RWMutex
	Routes             map[string]*url.URL
	Transport          *http.Transport
	Md5                string
	LokiUrl            string
	RateLimiterStorage storage.RateLimiterStorage
	CacherStorage      storage.CacherStorage
}
