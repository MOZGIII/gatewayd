package session

import (
	"gatewayd/backend/profile"
	"gatewayd/driver"
)

// Session stores internal session information.
// Session does not know it's token, only manager does.
type Session struct {
	profile *profile.Profile  // profile defines session settings
	driver  driver.Driver     // driver does the actual work
	params  map[string]string // params are passed on session creation

	donech chan struct{}
}

// New makes a new session struct.
func New(profile *profile.Profile, driver driver.Driver, params map[string]string) *Session {
	return &Session{
		profile, driver, params,
		make(chan struct{}),
	}
}

// Create makes new session by profile specification.
func Create(profile *profile.Profile, params map[string]string) (*Session, error) {
	driver, err := driver.CreateByName(profile.DriverName)
	if err != nil {
		return nil, err
	}
	session := New(profile, driver, params)
	if err := driver.Assign(session); err != nil {
		return nil, err
	}
	return session, nil
}

// Driver returns session driver.
func (s *Session) Driver() driver.Driver {
	return s.driver
}

// Profile returns session profile.
func (s *Session) Profile() *profile.Profile {
	return s.profile
}

// Params returns session params.
func (s *Session) Params() map[string]string {
	return s.params
}
