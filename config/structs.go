package config

type Endpoint struct {
	Host    string `yaml:"host"`
	Service string `yaml:"service"`
}

type Service struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

type AppConfig struct {
	Services             []Service  `yaml:"services"`
	Endpoints            []Endpoint `yaml:"endpoints"`
	Address              string     `yaml:"proxy.address"`
	ReloadConfigInterval int        `yaml:"config.reload_interval"`
	LokiUrl              string     `yaml:"monitoring.logging.loki"`
}
