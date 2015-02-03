package config

// ServerEndpoint is used by config to set
// addresses and SSL settings
type ServerEndpoint struct {
	Addr string

	SSLEnabled  bool
	SSLKeyFile  string
	SSLCertFile string
}
