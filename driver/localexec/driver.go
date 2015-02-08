package localexec

import (
	"fmt"
	"log"
	"net"
	"time"

	"gatewayd/driver"
	"gatewayd/driver/state"

	"gatewayd/backend/abstract"
)

type localExecDriver struct {
	session      abstract.Session
	state        state.Type
	stateChanged chan state.Type
}

// NewLocalExecDriver is a factory function to create new driver.
func NewLocalExecDriver() driver.Driver {
	return &localExecDriver{
		nil,
		state.Stopped,
		make(chan state.Type),
	}
}

func init() {
	driver.FactoryRegister("localexec", NewLocalExecDriver)
}

func (l *localExecDriver) Assign(session abstract.Session) error {
	if l.session != nil {
		fmt.Errorf("localexec: session already assigned")
	}

	l.session = session
	return nil
}

func (l *localExecDriver) Start() error {
	log.Println("localexec: start mock!")

	// Simulate session being started and stopped
	go func() {
		{
			t := 200 * time.Millisecond
			log.Printf("localexec: mock will boot up in %s", t)
			time.Sleep(t)
		}
		log.Printf("localexec: mock booted up")
		l.stateChanged <- state.Started

		{
			t := 5000 * time.Millisecond
			log.Printf("localexec: mock will shutdown in %s", t)
			time.Sleep(t)
		}
		log.Printf("localexec: mock shutdown")
		l.stateChanged <- state.Stopped
	}()

	return nil
}

func (l *localExecDriver) State() state.Type {
	return l.state
}

func (l *localExecDriver) StateChanged() <-chan state.Type {
	return l.stateChanged
}

func (l *localExecDriver) Terminate() error {
	log.Println("localexec: terminate mock!")
	return nil
}

func (l *localExecDriver) RemoteVNCConnection() (net.Conn, error) {
	log.Println("localexec: RemoteVNCConnection mock!")
	return net.Dial("tcp", "127.0.0.1:6900")
}
