package proxy

import (
	"log"
	"net/http"
	"net/url"
	"voltgate-proxy/config"
)

func NewProxyServer(config *config.AppConfig) *Server {
	proxy := &Server{
		routes:    make(map[string]*url.URL),
		transport: &http.Transport{},
	}

	for _, service := range config.Services {
		for _, endpoint := range config.Endpoints {
			if endpoint.Service == service.Name {
				parsedURL, err := url.Parse(service.Url)
				if err == nil {
					proxy.routes[endpoint.Host] = parsedURL
					log.Println("Mapping:", endpoint.Host, "->", service.Url)
				}
			}
		}
	}

	return proxy
}
