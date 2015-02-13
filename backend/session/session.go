package session

import (
	"gatewayd/backend/profile"
	"gatewayd/driver"
)

// Session stores internal session information.
// Session does not know it's token, only manager does.
type Session struct {
	profile *profile.Profile // profile defines session settings
	driver  driver.Driver    // driver does the actual work

	donech chan struct{}
}

// New makes a new session struct.
func New(profile *profile.Profile, driver driver.Driver) *Session {
	return &Session{
		profile, driver,
		make(chan struct{}),
	}
}

// Create makes new session by profile specification.
func Create(profile *profile.Profile) (*Session, error) {
	driver, err := driver.CreateByName(profile.DriverName)
	if err != nil {
		return nil, err
	}
	session := New(profile, driver)
	if err := driver.Assign(session); err != nil {
		return nil, err
	}
	return session, nil
}

// Driver returns session driver.
func (s *Session) Driver() driver.Driver {
	return s.driver
}
