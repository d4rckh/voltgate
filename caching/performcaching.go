package caching

import (
	"log"
	"net/http"
	"regexp"
	"time"
	"voltgate-proxy/config"
	"voltgate-proxy/proxy"
)

func PerformCaching(p *proxy.Server, cacheRules []config.CacheRule, r *http.Request, rwTrap *proxy.ResponseWriterTrap, doRequest func(rwTrap *proxy.ResponseWriterTrap)) {
	servedRequest := false

	for _, cacheRule := range cacheRules {
		matchedPath, _ := regexp.Match(cacheRule.Path, []byte(r.URL.Path))

		if !matchedPath {
			continue
		}

		status, header, body, exists := p.CacherStorage.GetRequest(r.Method, r.URL.Path)

		if !exists {
			log.Printf("doing request @ %s", r.URL.Path)

			doRequest(rwTrap)
			servedRequest = true

			p.CacherStorage.CacheRequest(r.Method, r.URL.Path, rwTrap.StatusCode, rwTrap.Header(), rwTrap.Body, time.Duration(cacheRule.Ttl)*time.Second)
			return
		}

		log.Printf("serving cached request @ %s", r.URL.Path)

		for key, values := range header {
			for _, value := range values {
				rwTrap.Header().Add(key, value)
			}
		}
		rwTrap.WriteHeader(status)
		rwTrap.Write(body)

		servedRequest = true
	}

	if !servedRequest {
		doRequest(rwTrap)
	}

}
