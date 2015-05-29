package global

import (
	"gatewayd/profile"
	"gatewayd/runner"
	"gatewayd/session"
)

// Store some stuff globally.
var (
	SessionRegistry = *session.NewRegistry()
	ProfileManager  = *profile.NewManager()
	Runner          = runner.NewRunner()
)
