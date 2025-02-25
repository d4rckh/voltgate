# voltgate proxy

### Example config.yaml
```yaml
proxy.address: ":80"

services:
  - name: prometheus
    url: http://localhost:9090

endpoints:
  - host: resedinta.dfourmusic.com
    service: prometheus
```