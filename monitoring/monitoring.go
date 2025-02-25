package monitoring

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func MonitorRequest(p *proxy.Server, request *http.Request, originalUrl *url.URL, code int, size int, duration time.Duration) {
	logMsg := fmt.Sprintf("[%s] -> [%s] -> %s %s (%d / %d bytes / %dms)",
		request.RemoteAddr, originalUrl.Host, request.Method, request.URL, code, size, duration.Milliseconds())

	log.Printf(logMsg)

	go sendToLoki(p, logMsg)
}
