package control

import (
	"gatewayd/backend/profile"
	"gatewayd/backend/session"
)

// GlobalContext stores cotext globally.
type GlobalContext struct {
	SessionManager session.Manager
	ProfileManager profile.Manager
}

// Do not export this so we can only access non-threadsafe methods from here
var context = GlobalContext{
	SessionManager: *session.NewManager(),
}

// FixmeContextExport  FIXME!!!
func FixmeContextExport() *GlobalContext {
	return &context
}
