package api

import (
	"log"
	"net/http"

	control "gatewayd/backend/control"
	"gatewayd/tunnel"
)

// This will soon move out from the api package...

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

	log.Printf("New tunnel connection with token %q", token)

	_, err := control.FixmeContextExport().SessionManager.SessionByToken(token)
	if err != nil {
		log.Printf("Blocked access attempt with invalid token %q", token)
		log.Println(err)
		http.Error(rw, "Access denied", 403)
		return
	}

	tunnel.Handler(rw, req)
}
