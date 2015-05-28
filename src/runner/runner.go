package runner

import (
	"errors"
	"log"
	"sync"
)

// Runner runs and manages runables.
type Runner interface {
	Run(runable Runable) error // executes passed Runable in a separate goroutine
	Go(run RunFunc) error      // executes passed RunFunc in a separate goroutine
	TerminateAll()             // send all maanged runables command to terminate
	Wait()                     // executes when all runables are finished execution
}

type runner struct {
	closech chan struct{}
	wg      sync.WaitGroup
}

func (r *runner) Run(runable Runable) error {
	return r.Go(runable.Run)
}

func (r *runner) Go(run RunFunc) error {
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
			log.Printf("runable: job finished %v", run)
		}()
		run(r.closech)
	}()
	log.Printf("runable: job started %v", run)
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
