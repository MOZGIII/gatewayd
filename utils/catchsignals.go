package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// CatchSignals catches signals
func CatchSignals(quit chan<- struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c

	timeout := 15 * time.Second

	log.Printf("Captured %v, exiting (force exit in %s)...", sig, timeout)

	go func() {
		time.Sleep(timeout)
		log.Printf("Passed %s, force exit...", timeout)
		os.Exit(1)
	}()

	quit <- struct{}{}
}
