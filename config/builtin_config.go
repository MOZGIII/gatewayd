package config

var defaultConfig ConfigStruct = ConfigStruct{
	PublicEndpoint: ServerEndpoint{
		Addr:        ":3000",
		SSLEnabled:  false,
		SSLKeyFile:  "public.key",
		SSLCertFile: "public.crt",
	},
	ServiceEndpoint: ServerEndpoint{
		Addr:        ":3001",
		SSLEnabled:  false,
		SSLKeyFile:  "service.key",
		SSLCertFile: "service.crt",
	},
}

// LoadBuiltinConfig loads default bundled config
func LoadBuiltinConfig() *ConfigStruct {
	return &defaultConfig
}
