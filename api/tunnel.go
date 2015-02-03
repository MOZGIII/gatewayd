package api

import (
	"log"
	"net/http"

	"gatewayd/backend"
	"gatewayd/tunnel"
)

// Tunnel is a handler that checks request for validity and then
// calls tunnel.Handler to trigger tunnel
// It's sort of authentication wrapper
func Tunnel(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(rw, "Method not allowed", 405)
		return
	}

	token := req.URL.Query().Get("token")
	if token == "" {
		log.Printf("Blocked access attempt without token")
		http.Error(rw, "No token passed", 400)
		return
	}

	log.Printf("New connection with token %q", token)

	ansc := make(chan *backend.Session)
	backend.Control.SyncChannel() <- func(sm *backend.SessionsManager) {
		ansc <- sm.SessionByToken(token)
	}
	session := <-ansc

	if session == nil {
		log.Printf("Blocked access attempt with invalid token %q", token)
		http.Error(rw, "Access denied", 403)
		return
	}

	tunnel.Handler(rw, req)
}
