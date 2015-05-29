package driver

import (
	"gatewayd/profile"
)

// Session is an interface that drivers use to
// working with sessions.
type Session interface {
	Profile() *profile.Profile
}
