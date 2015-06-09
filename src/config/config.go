package config

import (
	"encoding/json"
	"os"
)

// ServerEndpoint is used by config to set addresses
// and SSL settings.
type ServerEndpoint struct {
	Addr string `json:"addr"`

	SSLEnabled  bool   `json:"ssl_enabled"`
	SSLKeyFile  string `json:"ssl_keyfile"`
	SSLCertFile string `json:"ssl_certfile"`
}

// PolicyConfig is be used by config to load policies
// configuration.
type PolicyConfig struct {
	SessionManagementPolicyName string `json:"session_management"`
	TunnelAccessPolicyName      string `json:"tunnel_access"`
}

// Config is stored as JSON and used for app settings.
type Config struct {
	PublicEndpoint  ServerEndpoint `json:"public_endpoint"`
	ServiceEndpoint ServerEndpoint `json:"service_endpoint"`
	Policies        PolicyConfig   `json:"policies"`
}

// LoadConfig loads configuration from JSON file.
func LoadConfig(filename string) (*Config, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// MustLoadConfig calls LoadConfig and makes it a must-success.
func MustLoadConfig(filename string) (config *Config) {
	config, err := LoadConfig(filename)
	if err != nil {
		panic(err)
	}
	return
}
