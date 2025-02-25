package config

import (
	"log"
	"net/url"
	"time"
	"voltgate-proxy/proxy"
)

func BuildProxyFromConfig(proxy *proxy.Server, config *AppConfig, md5 string) {
	proxy.Mu.Lock()
	defer proxy.Mu.Unlock()

	proxy.Routes = make(map[string]*url.URL)

	for _, service := range config.Services {
		for _, endpoint := range config.Endpoints {
			if endpoint.Service == service.Name {
				parsedURL, err := url.Parse(service.Url)
				if err == nil {
					proxy.Routes[endpoint.Host] = parsedURL
					log.Println("Mapping:", endpoint.Host, "->", service.Url)
				}
			}
		}
	}

	log.Println("Successfully loaded", len(config.Endpoints), "endpoints and", len(config.Services), "services")

	proxy.Md5 = md5
}

func LoadConfig(proxy *proxy.Server, filename string) *AppConfig {
	config, md5, err := ReadConfig(filename)

	if err != nil {
		panic(err)
	}

	BuildProxyFromConfig(proxy, config, md5)

	return config
}

func ReloadConfig(proxyServer *proxy.Server, secondsInterval int, filename string) {
	ticker := time.NewTicker(time.Duration(secondsInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		config, md5, err := ReadConfig(filename)

		if err != nil {
			log.Printf("Error reloading config: %v", err)
			return
		}

		if md5 != proxyServer.Md5 {
			log.Println("Detected change in config, reloading routes")
			BuildProxyFromConfig(proxyServer, config, md5)
		}
	}
}
