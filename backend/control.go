package backend

import (
	"log"

	"gatewayd/utils"
)

// Do not export this so we can only access non-threadsafe methods from here
var sessionsManager = NewSessionsManager()

// ControlFunc definition is used when talking to ControlStorage
type ControlFunc func(sm *SessionsManager)

// Control is a global variable used to manage data in the app
var Control = ControlType{
	syncChan:  make(chan ControlFunc),
	quit:      make(chan utils.Empty), // buffered, to increase reliability
	debugChan: make(chan utils.Empty, 1),
}

// ControlType contains channels you should use to manage data
// in a thread-safe way
type ControlType struct {
	syncChan  chan ControlFunc
	quit      chan utils.Empty
	debugChan chan utils.Empty
}

// SyncChannel channel can be used to run function in a Control's goroutine,
// thus allowing for stnchronized race-free data access
func (c ControlType) SyncChannel() chan<- ControlFunc {
	return c.syncChan
}

// Quit can be used to gcaefully finish control loop
func (c ControlType) Quit() chan<- utils.Empty {
	return c.quit
}

// DebugChan makes control print some debug info upon being written to
func (c ControlType) DebugChan() chan<- utils.Empty {
	return c.debugChan
}

func (c ControlType) controlLoop() {
	for {
		select {
		case f := <-c.syncChan:
			f(sessionsManager)
		case <-c.quit:
			log.Println("Control loop quitting")
			return
		case <-c.debugChan:
			log.Println(sessionsManager)
		}
	}
}

// Run is a routine that runs control operations
func (c ControlType) Run() {
	log.Println("Control routine started")
	defer log.Println("Control routine finished")

	c.controlLoop()
}
