package management

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"voltgate-proxy/config"
)

func StartManagementServer(config *config.AppConfig) {
	log.Printf("Starting management server on %s", config.ManagementConfig.Address)
	if config.MonitoringAppConfig.PrometheusEnabled {
		log.Printf("Serving Prometheus metrics on %s/metrics", config.ManagementConfig.Address)
		http.Handle("/metrics", promhttp.Handler())
	}
	log.Fatal(http.ListenAndServe(config.ManagementConfig.Address, nil))
}
