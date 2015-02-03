package config

import (
	"encoding/json"
	"os"
)

// ConfigStruct is sotred as JSON and used to set some settings
type ConfigStruct struct {
	PublicEndpoint  ServerEndpoint `json:"public_endpoint"`
	ServiceEndpoint ServerEndpoint `json:"service_endpoint"`
}

// LoadConfig load configuration from json file
func LoadConfig(filename string) (*ConfigStruct, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)

	var config *ConfigStruct
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

// MustLoadConfig calls LoadConfig and makes it a must-success
func MustLoadConfig(filename string) (config *ConfigStruct) {
	config, err := LoadConfig(filename)
	if err != nil {
		panic(err)
	}
	return
}
