package sessionmanagement

import (
	"net/http"
)

// Policy represents an interface for
// session management policy implementations to look up to.
// This is called suring session creation. Implementaion has
// a chance to terminate up old/unused sessions that
// consume resources to allow new session to spawn, etc.
type Policy interface {
	SessionCreationHook(r *http.Request) error
}
