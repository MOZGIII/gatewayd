package session

import (
	"log"

	"gatewayd/pkg/tokenmap"
)

// Registry maps sessions and tokens to each other.
type Registry struct {
	sessionByToken *tokenmap.TokenMap
}

// NewRegistry creates new Registry
func NewRegistry() *Registry {
	return &Registry{tokenmap.New()}
}

// Manage registeres session immediately and manages session
// unregistration upon it's termination.
func (s *Registry) Manage(session *Session) (string, error) {
	token, err := s.Register(session)
	if err != nil {
		return "", err
	}
	log.Printf("sessionregistry: session %v registered with token %q", session, token)
	go s.unregisterOnTerminate(session)
	return token, err
}

// Register maps given session with some token and returns the token.
func (s *Registry) Register(session *Session) (string, error) {
	token, err := s.sessionByToken.InsertWithRandomKey(session)
	return token, err
}

// SessionByToken fetches session by it's token
func (s *Registry) SessionByToken(token string) (*Session, error) {
	val, err := s.sessionByToken.Get(token)
	if err != nil {
		return nil, err
	}
	return val.(*Session), nil
}

// Unregister removes session to token assosiation, making token obsolete.
func (s *Registry) Unregister(session *Session) error {
	return s.sessionByToken.RemoveByVal(session)
}

func (s *Registry) unregisterOnTerminate(session *Session) {
	log.Printf("sessionregistry: session termination watchdog for session %v started", session)
	defer log.Printf("sessionregistry: session termination watchdog for session %v finished", session)
	<-session.donech
	s.Unregister(session)
	log.Printf("sessionregistry: session %v unregistered", session)
}
