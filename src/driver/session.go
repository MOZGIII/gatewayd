package driver

import (
	"gatewayd/backend/profile"
)

// Session is an interface that drivers use to
// working with sessions.
type Session interface {
	Profile() *profile.Profile
}
