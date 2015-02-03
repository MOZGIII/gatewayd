package tunnel

import (
	"log"
	"net"
	"net/http"
)

// Handler is net.http compatile handler func that serves WebSocket
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("New connection from: %s", r.RemoteAddr)

	// Init external connection
	externalConn, err := net.Dial("tcp", "127.0.0.1:6900")
	if err != nil {
		log.Println("FAIL (net dial tcp):", err.Error())
		return
	}

	if err = Tunnel(w, r, externalConn); err != nil {
		log.Println("FAIL (tunnel):", err.Error())
		return
	}
}
