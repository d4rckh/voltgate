package proxy

import (
	"log"
	"net/http"
	"net/url"
)

func (p *Server) LogRequest(r *http.Request, u *url.URL) {
	log.Printf("Forwarding request from %s to %s", u.Host, r.URL)
}
