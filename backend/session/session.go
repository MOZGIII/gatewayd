package session

import (
	"log"

	// "gatewayd/backend/control"
	"gatewayd/backend/profile"
	"gatewayd/driver"
)

// Session stores internal session information.
// Session does not know it's token, only manager does.
type Session struct {
	profile *profile.Profile // profile defines session settings
	driver  driver.Driver    // driver does the actual work

	terminateChannel chan struct{}
}

// New creates a new session.
func New(profile *profile.Profile, driver driver.Driver) *Session {
	return &Session{
		profile, driver,
		make(chan struct{}),
	}
}

// Terminate can be called to terminate this session.
func (s *Session) Terminate() {
	close(s.terminateChannel)
}

func (s *Session) run() {
	select {
	case <-s.terminateChannel:
		log.Printf("Session %v is terminating", s)
		return
	}
}
