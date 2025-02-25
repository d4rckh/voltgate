package proxy

import (
	"net/http"
	"net/url"
	"sync"
	"voltgate-proxy/rate_limiting"
)

type Server struct {
	Mu                     sync.RWMutex
	Routes                 map[string]*url.URL
	Transport              *http.Transport
	Md5                    string
	LokiUrl                string
	EndpointRateLimitRules map[string][]rate_limiting.RateLimitRule
	ServicesRateLimitRules map[string][]rate_limiting.RateLimitRule
	RateLimiterStorage     rate_limiting.RateLimiterStorage
}
