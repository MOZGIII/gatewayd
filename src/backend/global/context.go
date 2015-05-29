package global

import (
	"gatewayd/backend/profile"
	"gatewayd/backend/session"
	"gatewayd/runner"
)

// Store some stuff globally.
var (
	SessionRegistry = *session.NewRegistry()
	ProfileManager  = *profile.NewManager()
	Runner          = runner.NewRunner()
)
