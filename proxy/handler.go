package proxy

import (
	"net/http"
	"net/http/httputil"
	"time"
)

func (p *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	targetURL, exists := p.Routes[r.Host]
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = p.transport

	originalURL := *r.URL
	originalURL.Host = r.Host

	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Host = targetURL.Host

	rwTrap := ResponseWriterTrap{ResponseWriter: w}

	starTime := time.Now()
	proxy.ServeHTTP(&rwTrap, r)
	duration := time.Since(starTime)

	p.MonitorRequest(r, &originalURL, rwTrap.StatusCode, rwTrap.ContentSize, duration)
}
