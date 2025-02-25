package rate_limiting

type RateLimitRule struct {
	Path             string `yaml:"path"`
	NumberOfRequests int    `yaml:"requests"`
	WindowSeconds    int    `yaml:"window"`
	Action           string `yaml:"action"`
	Method           string `yaml:"method"`
}
