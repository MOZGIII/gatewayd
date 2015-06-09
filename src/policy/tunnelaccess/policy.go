package tunnelaccess

import (
	"net/http"

	"gatewayd/session"
)

// Policy describes the interface to be implemented by all policies
// that want to hook into tunnel access management.
type Policy interface {
	TunnelAccessHook(s *session.Session, r *http.Request) error
}
