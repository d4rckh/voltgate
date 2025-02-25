package proxy

import (
	"net/http"
	"net/url"
	"sync"
)

type Server struct {
	mu        sync.RWMutex
	routes    map[string]*url.URL
	transport *http.Transport
}
