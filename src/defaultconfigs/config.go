package defaultconfigs

import (
	"gatewayd/config"
)

var defaultConfig = config.Config{
	PublicEndpoint: config.ServerEndpoint{
		Addr:        ":3000",
		SSLEnabled:  false,
		SSLKeyFile:  "public.key",
		SSLCertFile: "public.crt",
	},
	ServiceEndpoint: config.ServerEndpoint{
		Addr:        ":3001",
		SSLEnabled:  false,
		SSLKeyFile:  "service.key",
		SSLCertFile: "service.crt",
	},
}

// GetConfig loads default bundled config.
func GetConfig() *config.Config {
	return &defaultConfig
}
