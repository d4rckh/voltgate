package config

import (
	"crypto/md5"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadConfig(filename string) (*AppConfig, string, error) {
	var appConfig AppConfig

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}

	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		return nil, "", err
	}

	endpointsJson, endpointsJsonErr := json.Marshal(appConfig.Endpoints)
	servicesJson, servicesJsonErr := json.Marshal(appConfig.Services)

	if servicesJsonErr != nil || endpointsJsonErr != nil {
		return nil, "", err
	}

	sum := md5.Sum(append(endpointsJson, servicesJson...))

	return &appConfig, string(sum[:]), nil
}
