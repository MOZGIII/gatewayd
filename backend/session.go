package backend

import (
	"errors"
	"sync"

	"gatewayd/utils"
)

// Session stores internal session information
type Session struct {
	key   string `json:"key"` // passed by user, identifies session
	token string `json:"key"` // generated randomly and returned to user
}

// TODO: add interface here to block global instantiation, maybe

// Key returns session key
func (s Session) Key() string {
	return s.key
}

// Token returns session token
func (s Session) Token() string {
	return s.token
}

// Export generates object to be exposed to the APIs
func (s Session) Export() interface{} {
	return struct {
		Key   string `json:"key"`
		Token string `json:"token"`
	}{
		s.key,
		s.token,
	}
}

// SessionsManager controls sessions
// This class is not thread-safe, wrap it with your
// own concurrency protection
type SessionsManager struct {
	sessionByKey   map[string]*Session
	sessionByToken map[string]*Session

	mu sync.RWMutex
}

// NewSessionsManager asd
func NewSessionsManager() *SessionsManager {
	return &SessionsManager{
		sessionByKey:   make(map[string]*Session),
		sessionByToken: make(map[string]*Session),
	}
}

// CreateSession creates new session. maintaining hashes properly
func (s *SessionsManager) CreateSession(key string) (*Session, error) {
	if sess, ok := s.sessionByKey[key]; ok {
		return sess, errors.New("A session already exists for this key")
	}

	token := s.GetFreeToken()
	session := &Session{
		key:   key,
		token: token,
	}
	s.registerSession(session)

	return session, nil
}

// GetFreeToken generates you a free token, that is not currently in
// the SessionsManager
func (s *SessionsManager) GetFreeToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for {
		token := utils.RandStr(32)
		if _, ok := s.sessionByToken[token]; !ok {
			return token
		}
	}
}

func (s *SessionsManager) registerSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionByKey[session.key] = session
	s.sessionByToken[session.token] = session
}

func (s *SessionsManager) unregisterSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessionByKey, session.key)
	delete(s.sessionByToken, session.token)
}

// SessionByKey fetches session by it's key
func (s *SessionsManager) SessionByKey(key string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessionByKey[key]
}

// SessionByToken fetches session by it's token
func (s *SessionsManager) SessionByToken(token string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessionByToken[token]
}

// RemoveSession removes the session, unregistering it from the
// manager
func (s *SessionsManager) RemoveSession(session *Session) {
	// No locking here since we'll lock down the stack
	s.unregisterSession(session)
}
