# voltgate proxy

## features
- reverse proxy with support for multiple services and domains
- hot reloading of endpoints and services

### Example config.yaml
```yaml
proxy.address: ":80" # default: ":80"
proxy.config.reload_interval: 10 # reload services and endpoints every 10 seconds, default: do not reload

services:
  - name: prometheus
    url: http://localhost:9090

endpoints:
  - host: prometheus.host.com
    service: prometheus
```