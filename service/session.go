package service

import (
	// "os"
	"time"
)

type SessionInitParams struct {
	UnixUser  string
	UnixGroup string

	ScreenRes struct {
		W uint
		H uint
	}

	IdleTimeout time.Duration // nanoseconds
}

type Session struct {
	control chan int
}
