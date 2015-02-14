package localexec

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"gatewayd/driver"
	"gatewayd/driver/state"
)

type localExecDriver struct {
	session      driver.Session
	state        state.Type
	stateChanged chan state.Type

	cmd       *exec.Cmd
	localport int

	termch chan (chan error)
}

// NewLocalExecDriver is a factory function to create new driver.
func NewLocalExecDriver() driver.Driver {
	return &localExecDriver{
		nil,
		state.Init,
		make(chan state.Type),

		nil,
		0,

		make(chan (chan error)),
	}
}

func init() {
	driver.FactoryRegister("localexec", NewLocalExecDriver)
}

func (l *localExecDriver) Assign(session driver.Session) error {
	if l.session != nil {
		return fmt.Errorf("localexec: session already assigned to %v", l.session)
	}

	if l.state != state.Init {
		return fmt.Errorf("localexec: unable to assign session when driver is in %q state", l.state)
	}

	l.session = session
	l.initCommandFromSession()

	return nil
}

func (l *localExecDriver) Start() error {
	if l.state != state.Init {
		return fmt.Errorf("localexec: unable to start a session the is in %q state", l.state)
	}

	go l.run()
	return nil
}

func (l *localExecDriver) Terminate() error {
	log.Printf("localexec: terminate called")
	errch := make(chan error)
	l.termch <- errch
	err := <-errch
	return err
}

func (l *localExecDriver) State() state.Type {
	return l.state
}

func (l *localExecDriver) StateChanged() <-chan state.Type {
	return l.stateChanged
}

func (l *localExecDriver) RemoteVNCConnection() (net.Conn, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", l.localport)
	return net.Dial("tcp", addr)
}

func (l *localExecDriver) initCommandFromSession() error {
	if l.session == nil {
		return fmt.Errorf("localexec: unable to init command with nil session")
	}

	l.cmd = exec.Command("gatewayd-session-test")
	l.localport = 6500 // mock implementation

	return nil
}

// The internal goroutine.
func (l *localExecDriver) run() {
	log.Printf("localexec: goroutine started for %v", l)
	defer log.Printf("localexec: goroutine finished for %v", l)

	log.Println("localexec: session process starting")
	l.changeState(state.Starting)

	if err := l.cmd.Start(); err != nil {
		log.Println(err)
		log.Println("localexec: session process was unable to start")
		l.changeState(state.Stopped)
		return
	}

	log.Println("localexec: session process started")
	l.changeState(state.Started)

	// Wait for process to end in a separate routine.
	waitch := make(chan error)
	go func() { waitch <- l.cmd.Wait() }()

loop:
	for {
		select {
		// If monitored process exits, report error and finish
		// the loop. We no longer need to listen to events for it.
		case err := <-waitch:
			if err != nil {
				log.Println(err)
				log.Printf("localexec: session process terminated abnormally")
			}
			break loop

		// If we got a request on a termch channel, we should try
		// to kill process and report the result via provided channel.
		case errch := <-l.termch:
			errch <- l.terminate()
		}
	}

	log.Printf("localexec: session process stopped")
	l.changeState(state.Stopped)
}

// Should only be called from internal goroutine.
func (l *localExecDriver) terminate() error {
	// Can terminate if process is already stopped.
	if l.state == state.Stopped {
		return nil
	}

	// Trying to terminate the process that is not in the started state.
	if l.state != state.Started {
		log.Printf("localexec: termination of process at state %q ommited", l.state)
		return nil
	}

	// Kill triggers process exit, result will be propagated to Wait call.
	if err := l.cmd.Process.Kill(); err != nil {
		return err
	}

	// We are stopping now.
	l.changeState(state.Stopping)

	return nil
}

// Should only be called from internal goroutine.
func (l *localExecDriver) changeState(newState state.Type) {
	l.state = newState
	log.Printf("localexec: pushing new state %q", newState)
	l.stateChanged <- newState
	log.Printf("localexec: new state pushed %q", newState)
}
