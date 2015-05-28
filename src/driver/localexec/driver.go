package localexec

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"

	"gatewayd/driver"
	"gatewayd/driver/state"

	"github.com/jmoiron/jsonq"
)

type localExecDriver struct {
	session      driver.Session
	state        state.Type
	stateChanged chan state.Type

	cmd *exec.Cmd

	vncaddr   string       // stores VNC address
	vncaddrmu sync.RWMutex // guards access to VNC address

	termch chan (chan error)
}

// NewLocalExecDriver is a factory function to create new driver.
func NewLocalExecDriver() driver.Driver {
	return &localExecDriver{
		nil,
		state.Init,
		make(chan state.Type),

		nil, // cmd

		"",             // vncaddr
		sync.RWMutex{}, // vncaddrmu

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
	if err := l.initCommandFromSession(); err != nil {
		return err
	}

	return nil
}

func (l *localExecDriver) Start() error {
	if l.state != state.Init {
		return fmt.Errorf("localexec: unable to start session that is in %q state", l.state)
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
	l.vncaddrmu.RLock()
	defer l.vncaddrmu.RUnlock()

	if l.vncaddr == "" {
		return nil, fmt.Errorf("localexec: VNC connection address in not known yet")
	}

	return net.Dial("tcp", l.vncaddr)
}

func (l *localExecDriver) initCommandFromSession() error {
	if l.session == nil {
		return fmt.Errorf("localexec: unable to init command with nil session")
	}

	// Fetch all the values first
	jq := jsonq.NewQuery(l.session.Profile().GetParams())
	name, err := jq.String("command", "name")
	if err != nil {
		return err
	}
	args, err := jq.ArrayOfStrings("command", "args")
	if err != nil {
		return err
	}

	// Sanity checks
	if name == "" {
		return fmt.Errorf("localexec: profile param error: empty command")
	}

	l.cmd = exec.Command(name)
	if args != nil && len(args) > 0 {
		l.cmd.Args = append(l.cmd.Args, args...)
	}

	log.Printf("localexec: session assigned with command %v", l.cmd)
	return nil
}

// The internal goroutine.
func (l *localExecDriver) run() {
	log.Printf("localexec: goroutine started for %v", l)
	defer log.Printf("localexec: goroutine finished for %v", l)

	// Begin by switching state to starting.
	log.Println("localexec: session process starting")
	l.changeState(state.Starting)

	// Hook stdout.
	stdout, err := l.cmd.StdoutPipe()
	if err != nil {
		log.Printf("localexec: unable to hook stdout: %s", err)
		l.changeState(state.Stopped)
		return
	}

	// Start process.
	if err := l.cmd.Start(); err != nil {
		log.Printf("localexec: session process was unable to start: %s", err)
		l.changeState(state.Stopped)
		return
	}

	// Read VNC address from process' stdout.
	if err := l.readVNCAddress(stdout); err != nil {
		log.Printf("localexec: unable to get VNC address from child proccess: %s", err)
		l.changeState(state.Stopped)
		return
	}

	// Switch state to started now.
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

func (l *localExecDriver) readVNCAddress(r io.Reader) error {
	l.vncaddrmu.Lock()
	defer l.vncaddrmu.Unlock()

	if l.vncaddr != "" {
		return fmt.Errorf("localexec: VNC address already set to %q for %v", l.vncaddr, l)
	}

	var address struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r).Decode(&address); err != nil {
		return err
	}

	l.vncaddr = address.Address
	log.Printf("localexec: VNC address %q fetched for %v", l.vncaddr, l)

	return nil
}
