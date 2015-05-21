package tunnel

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
		Subprotocols:    []string{"binary"},
	}
)

// Tunnel upgrades HTTP connection to WebSocket and tunnels
func Tunnel(w http.ResponseWriter, r *http.Request, c net.Conn) error {
	log.Printf("New tunnel for %s to %s", r.RemoteAddr, c.RemoteAddr())

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	proxy := NewWebsocketProxy(ws, c)
	go proxy.doProxy() // do proxy handles connection closing
	return nil
}
