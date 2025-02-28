package monitoring

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"voltgate-proxy/proxy"
)

type LogEntryStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LogEntry struct {
	Streams []LogEntryStream `json:"streams"`
}

func MonitorRequest(p *proxy.Server, request *http.Request, originalUrl *url.URL, code int, size int, duration time.Duration, cached bool) {
	cachedMsg := ""

	if cached {
		cachedMsg = " [cached]"
	}

	logMsg := fmt.Sprintf("[%s] -> [%s] -> %s %s (%d / %d bytes / %dms)%s",
		request.RemoteAddr, originalUrl.Host, request.Method, request.URL, code, size, duration.Milliseconds(), cachedMsg)

	log.Printf("%s", logMsg)

	RequestCount.WithLabelValues(request.Method, originalUrl.Host, "", request.URL.Path, strconv.Itoa(code)).Inc()
	RequestDuration.WithLabelValues(request.Method, originalUrl.Host, "", request.URL.Path).Observe(float64(duration.Milliseconds()))
	ResponseSize.WithLabelValues(request.Method, originalUrl.Host, "", request.URL.Path).Observe(float64(size))

	go sendToLoki(p, logMsg, "info")
}

func MonitorBlockedRequest(p *proxy.Server, request *http.Request, count int, duration time.Duration) {
	logMsg := fmt.Sprintf("Blocked %s (%d requests in %f seconds) (Request URL: %s)",
		request.RemoteAddr, count, duration.Seconds(), request.URL)

	log.Printf("%s", logMsg)

	BlockedRequestCount.WithLabelValues(request.Method, request.Host, "", request.URL.Path).Inc()

	go sendToLoki(p, logMsg, "warn")
}
