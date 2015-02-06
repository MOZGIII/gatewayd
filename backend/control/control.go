package control

import (
	"log"
)

// Control is a global variable used to manage data in the app.
var Control = ControlType{
	make(chan struct{}),
	make(chan interface{}),
}

// ControlType are available for actors to send generic control events.
type ControlType struct {
	quit              chan struct{}
	sessionUnregister chan interface{} // only sessions should go here
}

// Quit can be used to gracefully finish control loop.
func (c ControlType) Quit() chan<- struct{} {
	return c.quit
}

// Run is a routine that runs control operations
func (c ControlType) Run() {
	log.Println("Control routine started")
	defer log.Println("Control routine finished")

	for {
		select {
		case <-c.quit:
			log.Println("Control loop: quit channel has spoken!")
			return
			// case s := <-c.sessionUnregister:
			// 	if err := context.SessionManager.Unregister(s); err != nil {
			// 		log.Println(err)
			// 	}
		}
	}
}
