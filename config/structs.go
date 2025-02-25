package config

type Endpoint struct {
	Host    string `yaml:"host"`
	Service string `yaml:"service"`
}

type Service struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

type MonitoringAppConfig struct {
	LokiUrl           string `yaml:"loki"`
	PrometheusEnabled bool   `yaml:"prometheus"`
}

type AppConfig struct {
	Services             []Service  `yaml:"services"`
	Endpoints            []Endpoint `yaml:"endpoints"`
	Address              string     `yaml:"proxy.address"`
	ManagementAddress    string     `yaml:"management.address"`
	ReloadConfigInterval int        `yaml:"config.reload_interval"`

	MonitoringAppConfig MonitoringAppConfig `yaml:"monitoring"`
}
