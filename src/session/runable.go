package session

import (
	"log"

	"gatewayd/driver/state"
)

// Run is the session management goroutine.
func (s *Session) Run(closech <-chan struct{}) {
	log.Printf("session: goroutine started for %v", s)
	defer log.Printf("session: goroutine finished for %v", s)

	// After routine ends...
	defer func() {
		// Trigger all external cleanup.
		close(s.donech)
	}()

	// Try starting driver.
	if err := s.driver.Start(); err != nil {
		log.Println(err)
		log.Printf("session: unable to start session %v", s)
		return
	}

	// Only take one from closech!
	closechOne := make(chan struct{})
	go func() {
		<-closech
		closechOne <- struct{}{}
	}()

	// Now react to state changes in a loop.
loop:
	for {
		select {
		case newState := <-s.driver.StateChanged():
			log.Printf("session: driver reported state change to %q (session %v)", newState, s)
			switch newState {
			case state.Stopped:
				break loop // driver stopped, we're done
			}
		case <-closechOne:
			// Runable graceful exit, order driver to terminate
			// and go on to wait for the state to become stopped.
			log.Printf("session: shutting down gracefully (session %v)", s)
			s.driver.Terminate()
		case change := <-s.tunnelBalanceCh:
			if change {
				s.tunnelsCount++
			} else {
				s.tunnelsCount--
			}
		}
	}
}
