package caching

import (
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"time"
	"voltgate-proxy/config"
	"voltgate-proxy/proxy"
)

func MakeCacheKey(r *http.Request, cacheRule config.CacheRule) string {
	cacheKey := r.URL.Host + r.URL.Path
	query := r.URL.Query()

	if cacheRule.Params == nil {
		// params not defined in config: add all query parameters
		if len(query) > 0 {
			cacheKey += "?" + query.Encode()
		}
	} else if len(cacheRule.Params) > 0 {
		// params defined but not empty: filter based on allowed parameters
		filteredQuery := url.Values{}
		for _, param := range cacheRule.Params {
			if values, exists := query[param]; exists {
				filteredQuery[param] = values
			}
		}
		if len(filteredQuery) > 0 {
			cacheKey += "?" + filteredQuery.Encode()
		}
	}
	return cacheKey
}

func PerformCaching(p *proxy.Server, cacheRules []config.CacheRule, r *http.Request, rwTrap *proxy.ResponseWriterTrap, doRequest func(rwTrap *proxy.ResponseWriterTrap)) bool {
	if p.CacherStorage == nil {
		doRequest(rwTrap)
		return false
	}

	for _, cacheRule := range cacheRules {
		cacheRule.Methods = []string{"GET"}

		matchedPath, _ := regexp.Match(cacheRule.Path, []byte(r.URL.Path))

		if !(matchedPath && (slices.Contains(cacheRule.Methods, r.Method) || len(cacheRule.Methods) == 0)) {
			continue
		}

		cacheKey := MakeCacheKey(r, cacheRule)

		status, header, body, exists := p.CacherStorage.GetRequest(r.Method, cacheKey)

		if !exists {
			//log.Printf("doing request @ %s", r.URL.Path)

			doRequest(rwTrap)

			p.CacherStorage.CacheRequest(r.Method, cacheKey, rwTrap.StatusCode, rwTrap.Header(), rwTrap.Body, time.Duration(cacheRule.Ttl)*time.Second)
			return false
		}

		//log.Printf("serving cached request @ %s", cacheKey)

		for key, values := range header {
			for _, value := range values {
				rwTrap.Header().Add(key, value)
			}
		}
		rwTrap.WriteHeader(status)
		rwTrap.Write(body)

		return true
	}

	doRequest(rwTrap)

	return false
}
