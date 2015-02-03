package backend

import "testing"

func TestSessionManager(t *testing.T) {
	go Control.Run()

	ansc := make(chan struct {
		session *Session
		err     error
	})
	Control.SyncChannel() <- func(sm *SessionsManager) {
		session, err := sm.CreateSession("mykey")
		ansc <- struct {
			session *Session
			err     error
		}{session, err}
	}
	ans := <-ansc

	if ans.err != nil {
		t.Error("Cannot create session", ans.err)
	}

	ansc2 := make(chan *Session)
	Control.SyncChannel() <- func(sm *SessionsManager) {
		ansc2 <- sm.SessionByKey("mykey")
	}
	gs := <-ansc2

	if gs == nil {
		t.Error("Cannot get just created session")
	}

	if gs != ans.session {
		t.Error("Session got in not the same one that was created!")
	}

	if ans.session.key != "mykey" {
		t.Error("Session has invalid key after all the tests!")
	}
}
