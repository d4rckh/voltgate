# Voltgate Proxy

## Features
- Reverse proxy supporting multiple services and domains
- Hot reloading of endpoints and services
- Log publishing to Loki

### Configuration Example (config.yaml)
```yaml
management:
  address: ":9999"
    # Default: ":9999"
proxy:
  address: ":80"
    # Default: ":80"

config:
  reload_interval: 10
    # Reloads monitoring, services, and endpoints every 10 seconds
    # Default: do not reload

monitoring:
  loki: http://localhost:3100/loki/api/v1/push
    # Publishes logs to Loki
    # Default: do not publish
  prometheus: true
    # Exposes /metrics on the management address
    # Default: false

services:
  - name: prometheus
    url: http://localhost:9090

endpoints:
  - host: prometheus.host.com
    service: prometheus
```

## Metrics Overview
Voltgate Proxy collects and exposes metrics to Prometheus for monitoring.

### `http_requests_total` (Counter)
Counts the total number of HTTP requests received.

- **Type:** Counter
- **Labels:**
    - `method`: HTTP method (e.g., GET, POST)
    - `host`: Original requested host
    - `target_service_name`: Proxied service name
    - `path`: Request path
    - `status`: HTTP response status code

#### Query Example:
```promql
rate(http_requests_total{method="GET", status="200"}[5m])
```
(Displays the rate of successful GET requests over 5 minutes.)

---
### `http_request_duration_seconds` (Histogram)
Measures the duration of HTTP requests in seconds.

- **Type:** Histogram
- **Buckets:** Default Prometheus latency buckets (`0.005s`, `0.01s`, `0.025s`, ...)
- **Labels:**
    - `method`: HTTP method
    - `host`: Requested host
    - `target_service_name`: Proxied service name
    - `path`: Request path

#### Query Example:
```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```
(Displays the 95th percentile of request durations over the last 5 minutes.)

---
### `http_response_size_bytes` (Histogram)
Tracks the size of HTTP responses in bytes.

- **Type:** Histogram
- **Buckets:** Exponential buckets (starting at 100 bytes, scaling by a factor of 2, up to 10 steps)
- **Labels:**
    - `method`: HTTP method
    - `host`: Requested host
    - `target_service_name`: Proxied service name
    - `path`: Request path

#### Query Example:
```promql
histogram_quantile(0.5, rate(http_response_size_bytes_bucket[5m]))
```
(Displays the median response size over the last 5 minutes.)

## Accessing Metrics
Prometheus metrics are available via an HTTP endpoint. Ensure the monitoring service is running and query the endpoint:

```
GET /metrics
```

### Prometheus Configuration
Add the following to `prometheus.yml` to enable scraping of the proxy service:
```yaml
scrape_configs:
  - job_name: "voltgate_proxy"
    static_configs:
      - targets: ["localhost:9999"] # Management address set in the configuration
```

