package runner

// Runable is used to implement structs with runable goroutines.
// The intention is to allow graceful termination of such runables.
type Runable interface {
	Run(closech <-chan struct{})
}

// RunFunc is the type for Runable.Run function.
type RunFunc func(closech <-chan struct{})
