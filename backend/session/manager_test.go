package session

import "testing"

func TestReadWrite(t *testing.T) {
	sm := NewSessionManager()

	session, err := sm.CreateSession("mykey")
	if err != nil {
		t.Error("Cannot create session", err)
	}

	gs := sm.SessionByKey("mykey")
	if gs == nil {
		t.Error("Cannot get just created session")
	}

	if gs != session {
		t.Error("Session got in not the same one that was created!")
	}
}
