package noop

import (
	"net/http"

	"gatewayd/policy/tunnelaccess"

	"gatewayd/session"
)

func init() {
	tunnelaccess.Factory.Register("noop", New)
}

// New creates a new noop policy.
func New() tunnelaccess.Policy {
	return &noop{}
}

type noop struct{}

func (n *noop) TunnelAccessHook(s *session.Session, r *http.Request) error {
	return nil
}
