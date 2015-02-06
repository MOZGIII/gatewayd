package driver

// Driver represnts an object capable of starting and managing
// a desktop session for us. It is responsible for strating up
// everything that is technically required for user to connect
// to: i.e. X and VNC server in LXC environment or KVM virtual
// machine with VNC display.
type Driver interface {
	StartSession() error     // starts session
	TerminateSession() error // terminates session (force kick user out)
	State() StateType
}

// StateType is for a driver's state machine type
type StateType int

const (
	// Stopped when session is stopped
	Stopped StateType = iota

	// Started when session is started
	Started StateType = iota
)
