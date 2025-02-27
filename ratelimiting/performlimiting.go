package ratelimiting

import (
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
	"voltgate-proxy/config"
	"voltgate-proxy/monitoring"
	"voltgate-proxy/proxy"
)

func PerformLimiting(p *proxy.Server, rules []config.RateLimitRule, request *http.Request) bool {
	clientAddr, _, err := net.SplitHostPort(request.RemoteAddr)

	if err != nil {

	}

	for _, rule := range rules {
		matchedPath, _ := regexp.Match(rule.Path, []byte(request.URL.Path))
		matchedRule := (matchedPath || rule.Path == "") && (rule.Method == request.Method || rule.Method == "*")

		if !matchedRule {
			continue
		}

		window := time.Duration(rule.WindowSeconds) * time.Second

		ruleKey := strings.Join([]string{clientAddr, rule.Path, rule.Method}, ",")

		requestCount, ttl := p.RateLimiterStorage.IncreaseRequestCount(ruleKey, window)

		//log.Printf("Found rule: requestCount(%d) ttl(%f)", requestCount, ttl.Seconds())

		if rule.NumberOfRequests < requestCount {
			monitoring.MonitorBlockedRequest(p, request, requestCount, window-ttl)

			return false
		}
	}

	return true
}
