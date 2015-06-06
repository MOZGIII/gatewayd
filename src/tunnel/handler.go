package tunnel

import (
	"log"
	"net/http"

	"gatewayd/global"
)

// Handler is dispatched for every tunnel connection. It validates
// the request and session status, makes a connection to the remote
// screen, upgrades HTTP connection to WebSocket and starts tunneling.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("tunnel: blocked access attempt without token")
		http.Error(w, "No token passed", 400)
		return
	}

	log.Printf("tunnel: new tunnel connection with token %q from %s", token, r.RemoteAddr)

	session, err := global.SessionRegistry.SessionByToken(token)
	if err != nil {
		log.Println(err)
		log.Printf("tunnel: blocked access attempt with invalid token %q", token)
		http.Error(w, "Access denied", 403)
		return
	}

	conn, err := session.Driver().RemoteVNCConnection()
	if err != nil {
		log.Println(err)
		log.Printf("tunnel: remote connection not ready for %q", token)
		http.Error(w, "Remote connection not ready", 503)
		return
	}

	log.Printf("New tunnel for %s to %s", r.RemoteAddr, conn.RemoteAddr())

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		log.Printf("tunnel: remote connection not ready for %q", token)
		// upgrader have already sent http error down the line
		return
	}
	// from now on, `w` is hijacked!

	tch := session.TunnelBalanceChannel()
	proxy := NewWebsocketProxy(ws, conn)
	go func() {
		defer func() {
			tch <- false // decrement tunnels counter
		}()
		proxy.doProxy() // do proxy handles connection closing
	}()
	tch <- true // increment tunnels counter
}
