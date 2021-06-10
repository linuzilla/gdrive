package interrupt_service

import (
	"os"
)

type InterruptService interface {
	Channel() chan os.Signal
	Exec(goroutine func(errChan chan error)) (bool, error)
	ExecWithChannel(errChan chan error, goroutine func(errChan chan error), waitRoutine func()) (bool, error)
	ClearInterrupt()
}
