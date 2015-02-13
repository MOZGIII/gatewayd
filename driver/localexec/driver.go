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
}

// NewLocalExecDriver is a factory function to create new driver.
func NewLocalExecDriver() driver.Driver {
	return &localExecDriver{
		nil,
		state.Stopped,
		make(chan state.Type, 1),

		nil,
		0,
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
	l.localport = 6500

	return nil
}

func (l *localExecDriver) Start() error {
	if err := l.cmd.Start(); err != nil {
		l.stateChanged <- state.Stopped
		return err
	}
	log.Println("localexec: session process started")

	go func() {
		if err := l.cmd.Wait(); err != nil {
			log.Println(err)
			log.Printf("localexec: session process terminated abnormally")
		}

		log.Printf("localexec: session process stopped")
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
	if l.cmd.Process == nil {
		// return fmt.Errorf("localexec: unable to kill, proccess in nil")
		return nil // proccess already dead
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
