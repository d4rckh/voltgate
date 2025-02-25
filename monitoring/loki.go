package monitoring

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"voltgate-proxy/proxy"
)

func sendToLoki(p *proxy.Server, logMsg string) {
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
