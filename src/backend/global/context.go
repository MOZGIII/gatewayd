package global

import (
	"gatewayd/backend/profile"
	"gatewayd/backend/session"
	"gatewayd/runner"
)

// Store some stuff globally.
var (
	SessionManager = *session.NewManager()
	ProfileManager = *profile.NewManager()
	Runner         = runner.NewRunner()
)
