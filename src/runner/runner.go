package runner

import (
	"errors"
	"log"
	"sync"
)

// Runner runs and manages runables.
type Runner interface {
	Go(runable Runable) error // executes runable in a separate goroutine
	TerminateAll()            // send all amanged runables signal to terminate
	Wait()                    // executes when all runables are finished execution
}

type runner struct {
	closech chan struct{}
	wg      sync.WaitGroup
}

func (r *runner) Go(runable Runable) error {
	select {
	case <-r.closech:
		// Got value on `closech`, which
		// means it's closed.
		// Block running new job.
		return errors.New("runner: running new jobs in not allowed for terminated runner")
	default:
		// Channel `closech` is blocked,
		// proceed with running.
	}

	r.wg.Add(1)
	go func() {
		defer func() {
			r.wg.Done()
			log.Printf("runable: job finished %v", runable)
		}()
		runable.Run(r.closech)
	}()
	log.Printf("runable: job started %v", runable)
	return nil
}

func (r *runner) TerminateAll() {
	close(r.closech)
}

func (r *runner) Wait() {
	r.wg.Wait()
}

// NewRunner is used to create and initialize new runners.
func NewRunner() Runner {
	return &runner{
		make(chan struct{}),
		sync.WaitGroup{},
	}
}
