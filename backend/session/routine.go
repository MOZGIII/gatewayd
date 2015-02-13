package session

import (
	"log"

	"gatewayd/driver/state"
)

// Run is the session management goroutine.
func (s *Session) Run() {
	log.Printf("session: goroutine started for %v", s)
	defer log.Printf("session: goroutine finished for %v", s)

	// After routine ends
	defer func() {
		// Call all internal cleanup
		if err := s.routineCleanup(); err != nil {
			log.Println(err)
			log.Printf("session: abnormal session termination of %v", s)
			panic(err) // crash everything on abnormal session termination!
		}

		// Trigger all external cleanup
		close(s.terminate)
	}()

	// Before routine starts
	if err := s.routineInit(); err != nil {
		log.Println(err)
		log.Printf("session: unable to start session %v", s)
		return
	}

	if err := s.routineMainloop(); err != nil {
		log.Println(err)
		log.Printf("session: session execution error at %v", s)
		return
	}
}

func (s *Session) routineInit() error {
	// Just simply start driver
	return s.driver.Start()
}

func (s *Session) routineCleanup() error {
	// We need to let driver cleanup
	return s.driver.Terminate()
}

func (s *Session) routineMainloop() error {
	for {
		select {
		// something happened on driver's side
		case newState := <-s.driver.StateChanged():
			log.Printf("session: driver reported state change to %q (session %v)", newState, s)
			switch newState {
			case state.Stopped:
				return nil // driver stopped, we're done
			}

		// we are to terminate this session due to internal command
		case <-s.terminate:
			return nil // do all the cleanup in routineCleanup func
		}
	}
}
