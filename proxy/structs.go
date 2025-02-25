package proxy

import (
	"net/http"
	"net/url"
	"sync"
)

type Server struct {
	Mu        sync.RWMutex
	Routes    map[string]*url.URL
	transport *http.Transport
	Md5       string
}
