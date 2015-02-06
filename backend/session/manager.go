package session

import (
	"gatewayd/pkg/tokenmap"
)

// Manager maps sessions and tokens to each other.
type Manager struct {
	sessionByToken *tokenmap.TokenMap
}

// NewManager creates new Manager
func NewManager() *Manager {
	return &Manager{tokenmap.New()}
}

// Register maps given session with some token and returns the token.
func (s *Manager) Register(session *Session) (string, error) {
	token, err := s.sessionByToken.InsertWithRandomKey(session)
	return token, err
}

// SessionByToken fetches session by it's token
func (s *Manager) SessionByToken(token string) (*Session, error) {
	val, err := s.sessionByToken.Get(token)
	if err != nil {
		return nil, err
	}
	return val.(*Session), nil
}

// Unregister removes session to token assosiation, making token obsolete.
func (s *Manager) Unregister(session *Session) error {
	return s.sessionByToken.RemoveByVal(session)
}
