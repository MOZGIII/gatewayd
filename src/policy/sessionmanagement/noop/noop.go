package noop

import (
	"net/http"

	"gatewayd/policy/sessionmanagement"
)

func init() {
	sessionmanagement.Factory.Register("noop", New)
}

// New creates a new noop policy.
func New() sessionmanagement.Policy {
	return &noop{}
}

type noop struct{}

func (n *noop) SessionCreationHook(r *http.Request) error {
	return nil
}
