package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func LoadConfig(filename string) (*AppConfig, error) {
	var appConfig AppConfig

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Loaded", len(appConfig.Endpoints), "endpoints and", len(appConfig.Services), "services")
	return &appConfig, nil
}
