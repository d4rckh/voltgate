package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "host", "target_service_name", "path", "status"},
	)

	BlockedRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_blocked_requests_total",
			Help: "Total number of blocked HTTP requests",
		},
		[]string{"method", "host", "target_service_name", "path"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "host", "target_service_name", "path"},
	)

	ResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "host", "target_service_name", "path"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestCount, RequestDuration, ResponseSize, BlockedRequestCount)
}
