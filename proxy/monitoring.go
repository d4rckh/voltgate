package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LogEntryStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LogEntry struct {
	Streams []LogEntryStream `json:"streams"`
}

func (p *Server) MonitorRequest(request *http.Request, originalUrl *url.URL, code int, size int, duration time.Duration) {
	logMsg := fmt.Sprintf("[%s] -> [%s] -> %s %s (%d / %d bytes / %dms)",
		request.RemoteAddr, originalUrl.Host, request.Method, request.URL, code, size, duration.Milliseconds())

	log.Printf(logMsg)

	go p.sendToLoki(logMsg)
}

func (p *Server) sendToLoki(logMsg string) {
	p.Mu.RLock()
	lokiUrl := p.LokiUrl
	p.Mu.RUnlock()

	logEntry := LogEntry{
		Streams: []LogEntryStream{ // Streams should be a slice
			{
				Stream: map[string]string{
					"job": "voltgate-server",
				},
				Values: [][]string{
					{strconv.FormatInt(time.Now().UnixNano(), 10), logMsg}, // Loki expects a timestamp and log message
				},
			},
		},
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	req, err := http.NewRequest("POST", lokiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf(err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("Failed to publish logs to Loki, received status code %d", resp.StatusCode)
	}
}
