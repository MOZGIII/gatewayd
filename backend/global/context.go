package global

import (
	"gatewayd/backend/profile"
	"gatewayd/backend/session"
)

// Store some stuff globally.
var (
	SessionManager = *session.NewManager()
	ProfileManager = *profile.NewManager()
)
