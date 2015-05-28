package driver

import (
	"net"

	"gatewayd/driver/state"
)

// Driver represnts an object capable of starting and managing
// a desktop session for us. It is responsible for strating up
// everything that is technically required for user to connect
// to: i.e. X and VNC server in LXC environment or KVM virtual
// machine with VNC display.
type Driver interface {
	Assign(session Session) error // assigns session to a driver

	// Actions
	Start() error            // starts session (probably should run driver's goroutine)
	Terminate() <-chan error // terminates session (force stops the session)

	// State machine
	State() state.Type               // the actual state of the driver for this session
	StateChanged() <-chan state.Type // channel to get the updates from the driver

	RemoteVNCConnection() (net.Conn, error) // returns a connection to tunnel to client
}
