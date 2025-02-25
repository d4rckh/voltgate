package proxy

import (
	"net/http"
	"net/url"
)

func NewProxyServer() *Server {
	proxy := &Server{
		Routes:    make(map[string]*url.URL),
		Transport: &http.Transport{},
	}

	return proxy
}
