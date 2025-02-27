package config

import (
	"voltgate-proxy/storage"
)

type RateLimitConfig struct {
	Rules []RateLimitRule `yaml:"rules"`
}

type AppRateLimitRules struct {
	EndpointRateLimitRules map[string][]RateLimitRule
	ServicesRateLimitRules map[string][]RateLimitRule
}

type Endpoint struct {
	Host            string          `yaml:"host"`
	Service         string          `yaml:"service"`
	RateLimitConfig RateLimitConfig `yaml:"rate_limit"`
}

type Service struct {
	Url             string          `yaml:"url"`
	Name            string          `yaml:"name"`
	RateLimitConfig RateLimitConfig `yaml:"rate_limit"`
}

type MonitoringAppConfig struct {
	LokiUrl           string `yaml:"loki"`
	PrometheusEnabled bool   `yaml:"prometheus"`
}

type RateLimitRule struct {
	Path             string `yaml:"path"`
	NumberOfRequests int    `yaml:"requests"`
	WindowSeconds    int    `yaml:"window"`
	Action           string `yaml:"action"`
	Method           string `yaml:"method"`
}

type RateLimitAppConfig struct {
	Storage string `yaml:"storage"`
}

type StorageAppConfig struct {
	Redis storage.RedisAppConfig `yaml:"redis"`
}

type AppConfig struct {
	Services             []Service  `yaml:"services"`
	Endpoints            []Endpoint `yaml:"endpoints"`
	Address              string     `yaml:"proxy.address"`
	ManagementAddress    string     `yaml:"management.address"`
	ReloadConfigInterval int        `yaml:"config.reload_interval"`

	Storage StorageAppConfig `yaml:"storage"`

	MonitoringAppConfig MonitoringAppConfig `yaml:"monitoring"`
	RateLimitAppConfig  RateLimitAppConfig  `yaml:"rate_limit"`
}
