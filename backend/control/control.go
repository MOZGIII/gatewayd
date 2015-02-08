package control

import (
	"log"
)

// Type are available for actors to send generic control events.
type Type struct {
	quit chan struct{}
}

// Control is a global variable used to manage data in the app.
var Control = Type{
	make(chan struct{}),
}

// Quit can be used to gracefully finish control loop.
func (c *Type) Quit() chan<- struct{} {
	return c.quit
}

// Run is a routine that runs control operations
func (c *Type) Run() {
	log.Println("control: goroutine started")
	defer log.Println("control: goroutine finished")

	for {
		select {
		case <-c.quit:
			log.Println("control: quit channel has spoken, shutting down")
			// TODO: terminate sessions here and wait for them to cleanup
			return
		}
	}
}
