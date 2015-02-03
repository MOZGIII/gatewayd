package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// CatchSignals catches signals
func CatchSignals(quit chan<- Empty) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c

	log.Printf("captured %v, exiting (force exit in 15 seconds)...", sig)

	go func() {
		time.Sleep(15 * time.Second)
		log.Printf("passed 15 seconds, force exit...")
		os.Exit(1)
	}()

	quit <- NewEmpty()
}
