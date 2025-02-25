package rate_limiting

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type RateLimiterStorage interface {
	IncreaseRequestCount(key string)
	GetRequestCount(key string) (int, time.Duration)
	ResetRequestCount(key string)
}

func PerformLimiting(storage RateLimiterStorage, rules []RateLimitRule, request *http.Request) bool {
	clientAddr := strings.Split(request.RemoteAddr, ":")[0]

	for _, rule := range rules {
		matchedPath, _ := regexp.Match(rule.Path, []byte(request.URL.Path))
		matchedRule := (matchedPath || rule.Path == "") && (rule.Method == request.Method || rule.Method == "*")

		if !matchedRule {
			continue
		}

		ruleKey := strings.Join([]string{clientAddr, rule.Path, rule.Method}, ",")
		requestCount, lastReset := storage.GetRequestCount(ruleKey)

		if lastReset > time.Duration(rule.WindowSeconds)*time.Second {
			storage.ResetRequestCount(clientAddr)
			storage.IncreaseRequestCount(clientAddr)
			return true
		}

		storage.IncreaseRequestCount(clientAddr)

		if rule.NumberOfRequests < requestCount {
			log.Printf("Blocked %s (%d requests in %fs) rule: %s '%s'", clientAddr, requestCount, lastReset.Seconds(), rule.Method, rule.Path)
			return false
		}
	}

	return true
}
