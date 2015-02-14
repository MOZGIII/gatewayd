package state

import (
	"fmt"
)

// Type is for a driver's state machine type
type Type int

const (
	// Init type is assigned by default for all new sessions
	// before any event happends on them.
	Init Type = iota

	// Starting is a state for when session was told to start
	// but the proccess of starting it is not yet finished.
	Starting

	// Started state is for when session had finished starting and
	// currently working.
	Started

	// Stopping is for when started session was ordered to stop but
	// the stopping proccess did not finish yet.
	Stopping

	// Stopped is a state for when session is stopped and
	// no longer working.
	Stopped
)

func (t Type) String() string {
	switch t {
	case Init:
		return "init"
	case Starting:
		return "starting"
	case Started:
		return "started"
	case Stopping:
		return "stopping"
	case Stopped:
		return "stopped"
	default:
		return fmt.Sprintf("undefined (%v)", int(t))
	}
}
