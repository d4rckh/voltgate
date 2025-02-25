package handler

import (
	"net/http"
	"net/http/httputil"
	"time"
	"voltgate-proxy/monitoring"
	"voltgate-proxy/proxy"
	"voltgate-proxy/rate_limiting"
)

func HandleRequest(p *proxy.Server, w http.ResponseWriter, r *http.Request) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	targetURL, exists := p.Routes[r.Host]
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	rateLimitingRules := p.EndpointRateLimitRules[r.Host]

	if !rate_limiting.PerformLimiting(p.RateLimiterStorage, rateLimitingRules, r) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.Transport = p.Transport

	originalURL := *r.URL
	originalURL.Host = r.Host

	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Host = targetURL.Host

	r.Header.Set("X-Forwarded-Host", originalURL.Host)
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)

	rwTrap := proxy.ResponseWriterTrap{ResponseWriter: w}

	startTime := time.Now()
	reverseProxy.ServeHTTP(&rwTrap, r)
	duration := time.Since(startTime)

	monitoring.MonitorRequest(p, r, &originalURL, rwTrap.StatusCode, rwTrap.ContentSize, duration)
}
