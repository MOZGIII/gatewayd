package session

import (
	"gatewayd/driver"
	"gatewayd/profile"
)

// Session stores internal session information.
// Session does not know it's token, only registry does.
type Session struct {
	profile *profile.Profile  // profile defines session settings
	driver  driver.Driver     // driver does the actual work
	params  map[string]string // params are passed on session creation

	tunnelsCount    uint32    // total amount of tunnels currently using this session
	tunnelBalanceCh chan bool // channel to keep tunnels balance with (true = +1, false = -1)

	donech chan struct{} // used by session registry to unregister session
}

// New makes a new session struct.
func New(profile *profile.Profile, driver driver.Driver, params map[string]string) *Session {
	return &Session{
		profile, driver, params,
		0, make(chan bool),
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

// TunnelsCount returns current tunnels count for this session.
func (s *Session) TunnelsCount() uint32 {
	return s.tunnelsCount
}

// TunnelBalanceChannel returns channel to be used to keep tunnels balance.
// Bool type, so true = +1, false = -1.
func (s *Session) TunnelBalanceChannel() chan<- bool {
	return s.tunnelBalanceCh
}
