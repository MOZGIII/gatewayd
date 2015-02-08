package state

import (
	"fmt"
)

// Type is for a driver's state machine type
type Type int

const (
	// Stopped when session is not running on host
	Stopped Type = iota

	// Started when session is is running on host
	Started Type = iota
)

func (t Type) String() string {
	switch t {
	case Started:
		return "started"
	case Stopped:
		return "stopped"
	default:
		return fmt.Sprintf("undefined (%v)", int(t))
	}
}
