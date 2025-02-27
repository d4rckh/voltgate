package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"voltgate-proxy/caching"
	"voltgate-proxy/config"
	"voltgate-proxy/monitoring"
	"voltgate-proxy/proxy"
	"voltgate-proxy/ratelimiting"
)

func ForwardToHost(p *proxy.Server, w *proxy.ResponseWriterTrap, targetURL *url.URL, r *http.Request) {
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.Transport = p.Transport

	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Host = targetURL.Host

	reverseProxy.ServeHTTP(w, r)
}

func HandleRequest(p *proxy.Server, rateLimitRules *config.AppRateLimitRules, cacheRules *config.AppCacheRules, w http.ResponseWriter, r *http.Request) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	targetURL, exists := p.Routes[r.Host]
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	rateLimitingRules := rateLimitRules.EndpointRateLimitRules[r.Host]
	cRules := cacheRules.EndpointCacheRules[r.Host]

	if !ratelimiting.PerformLimiting(p, rateLimitingRules, r) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	originalURL := *r.URL
	originalURL.Host = r.Host

	r.Header.Set("X-Forwarded-Host", originalURL.Host)
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)

	rwTrap := proxy.ResponseWriterTrap{ResponseWriter: w}

	startTime := time.Now()

	caching.PerformCaching(p, cRules, r, &rwTrap, func(rw *proxy.ResponseWriterTrap) {
		ForwardToHost(p, rw, targetURL, r)
	})

	duration := time.Since(startTime)

	monitoring.MonitorRequest(p, r, &originalURL, rwTrap.StatusCode, rwTrap.ContentSize, duration)
}
