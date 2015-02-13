package localexec

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"sync"

	"gatewayd/driver"
	"gatewayd/driver/state"
)

type localExecDriver struct {
	session      driver.Session
	state        state.Type
	stateChanged chan state.Type

	cmd       *exec.Cmd
	localport int

	mu sync.Mutex
}

// NewLocalExecDriver is a factory function to create new driver.
func NewLocalExecDriver() driver.Driver {
	return &localExecDriver{
		nil,
		state.Stopped,
		make(chan state.Type, 1),

		nil,
		0,

		sync.Mutex{},
	}
}

func init() {
	driver.FactoryRegister("localexec", NewLocalExecDriver)
}

func (l *localExecDriver) Assign(session driver.Session) error {
	if l.session != nil {
		return fmt.Errorf("localexec: session already assigned")
	}

	l.session = session
	l.initCommandFromSession()

	return nil
}

func (l *localExecDriver) initCommandFromSession() error {
	if l.session == nil {
		return fmt.Errorf("localexec: unable to init command with nil session")
	}

	l.cmd = exec.Command("gatewayd-session-test")
	l.localport = 6500 // mock implementation

	return nil
}

func (l *localExecDriver) Start() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.cmd.Start(); err != nil {
		l.changeState(state.Stopped)
		return err
	}
	log.Println("localexec: session process started")
	l.changeState(state.Started)

	go func() {
		// Wait for proccess to end.
		if err := l.cmd.Wait(); err != nil {
			log.Println(err)
			log.Printf("localexec: session process terminated abnormally")
		}

		// Lock further operations until done.
		l.mu.Lock()
		defer l.mu.Unlock()

		// Update the state.
		log.Printf("localexec: session process stopped")
		l.changeState(state.Stopped)
	}()

	return nil
}

func (l *localExecDriver) changeState(newState state.Type) {
	l.state = newState
	l.stateChanged <- newState
}

func (l *localExecDriver) State() state.Type {
	return l.state
}

func (l *localExecDriver) StateChanged() <-chan state.Type {
	return l.stateChanged
}

func (l *localExecDriver) Terminate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	log.Printf("localexec: terminate called")

	// Can terminate if process is already stopped
	if l.state == state.Stopped {
		return nil
	}

	// Kill triggers proccess exit, result will be propagated to Wait call.
	if err := l.cmd.Process.Kill(); err != nil {
		return err
	}
	return nil
}

func (l *localExecDriver) RemoteVNCConnection() (net.Conn, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", l.localport)
	return net.Dial("tcp", addr)
}
