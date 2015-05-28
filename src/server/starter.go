package server

import (
	"log"
	"net/http"

	"gatewayd/config"
)

// Start takes an handler and binds it to the endpoint
// The result is equal to call of ListenAndServe or ListenAndServeTLS
func Start(e config.ServerEndpoint, handler http.Handler) error {
	server := &http.Server{Addr: e.Addr, Handler: handler}

	if e.SSLEnabled {
		return server.ListenAndServeTLS(e.SSLKeyFile, e.SSLCertFile)
	}
	return server.ListenAndServe()
}

// RunAll starts servers for both of the endpoints,
// each in separate goroutines.
func RunAll(c *config.Config) {
	go func() {
		if err := Start(c.PublicEndpoint, NewPublicHander()); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := Start(c.ServiceEndpoint, NewServiceHander()); err != nil {
			log.Fatal(err)
		}
	}()
}
