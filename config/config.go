package config

import (
	"bytes"
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
	monitoringJson, monitoringJsonErr := json.Marshal(appConfig.MonitoringAppConfig.LokiUrl)

	if servicesJsonErr != nil || endpointsJsonErr != nil || monitoringJsonErr != nil {
		return nil, "", err
	}

	sum := md5.Sum(bytes.Join(
		[][]byte{endpointsJson, servicesJson, monitoringJson},
		[]byte{},
	))

	return &appConfig, string(sum[:]), nil
}
