package interrupt_service

import (
	"fmt"
	"os"
	"os/signal"
)

type interruptServiceImpl struct {
	signalChan chan os.Signal
}

func New() InterruptService {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	return &interruptServiceImpl{
		signalChan: interruptChannel,
	}
}

func (service *interruptServiceImpl) Channel() chan os.Signal {
	return service.signalChan
}

func (service *interruptServiceImpl) Exec(goroutine func(errChan chan error)) (bool, error) {
	service.ClearInterrupt()

	errChan := make(chan error, 1)

	defer close(errChan)

	return service.ExecWithChannel(errChan, goroutine, nil)
}

func (service *interruptServiceImpl) ExecWithChannel(errChan chan error, goroutine func(errChan chan error), waitRoutine func()) (bool, error) {
	service.ClearInterrupt()

	go goroutine(errChan)

	select {
	case err := <-errChan:
		return false, err
	case <-service.signalChan:
		fmt.Fprint(os.Stderr, "\n\n**** Break ****\n\n")
		return true, nil
	}
}

func (service *interruptServiceImpl) ClearInterrupt() {
	for {
		select {
		case <-service.signalChan:
			continue
		default:
			return
		}
	}
}
