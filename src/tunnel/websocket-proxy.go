package tunnel

import (
	"log"
	"net"

	"github.com/gorilla/websocket"
)

// WebsocketProxy stores a pair of WebSocket connection and
// a regular connection
type WebsocketProxy struct {
	ws    *websocket.Conn
	other net.Conn
}

// NewWebsocketProxy creates a new WebsocketProxy
func NewWebsocketProxy(ws *websocket.Conn, other net.Conn) *WebsocketProxy {
	p := WebsocketProxy{ws, other}
	return &p
}

func (p *WebsocketProxy) doProxy() {
	defer p.other.Close()
	defer p.ws.Close()

	done := make(chan bool)

	go func() { p.wsToOther(); done <- true }()
	go func() { p.otherToWs(); done <- true }()

	<-done // deferred stuff gets called here
}

func (p *WebsocketProxy) otherToWs() error {
	log.Println("VNC to WS goroutine started")
	defer func() {
		log.Println("VNC to WS goroutine ended")
	}()

	buffer := make([]byte, 1024)

	for {
		n, err := p.other.Read(buffer)
		if err != nil {
			log.Println("Error reading from VNC:", err.Error())
			return err
		}

		err = p.ws.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			log.Println("Error writing to WS:", err.Error())
			return err
		}
	}

}

func (p *WebsocketProxy) wsToOther() error {
	log.Println("WS to VNC goroutine started")
	defer func() {
		log.Println("WS to VNC goroutine ended")
	}()

	for {
		_, data, err := p.ws.ReadMessage()
		if err != nil {
			log.Println("Error reading from WS:", err.Error())
			return err
		}

		_, err = p.other.Write(data)
		if err != nil {
			log.Println("Error writing to VNC:", err.Error())
			return err
		}
	}
}
