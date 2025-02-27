package management

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"voltgate-proxy/config"
)

func StartManagementServer(config *config.AppConfig) {
	log.Printf("Starting management server on %s", config.ManagementAddress)
	if config.MonitoringAppConfig.PrometheusEnabled {
		log.Printf("Serving Prometheus metrics on %s/metrics", config.ManagementAddress)
		http.Handle("/metrics", promhttp.Handler())
	}
	log.Fatal(http.ListenAndServe(config.ManagementAddress, nil))
}
