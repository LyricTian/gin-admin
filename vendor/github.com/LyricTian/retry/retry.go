package retry

import (
	"errors"
	"time"
)

var (
	// ErrMaxRetries exceeded retry limit
	ErrMaxRetries = errors.New("exceeded retry limit")
)

type (
	// Func try to execute the function
	Func func() error

	// Sleep sleep callback function
	Sleep func(i int) time.Duration

	// Trier try to execute the interface
	Trier interface {
		Try() error
	}
)

// DoFunc try to execute a function,
// specify the number of attempts, and the execution interval
func DoFunc(retries int, fn Func, sleeps ...Sleep) error {
	if retries <= 0 {
		return nil
	}

	total := retries

LBBEGIN:

	if err := fn(); err != nil {
		retries--
		if retries == 0 {
			return ErrMaxRetries
		}

		if len(sleeps) > 0 {
			if d := sleeps[0](total - retries); d > 0 {
				time.Sleep(d)
			}
		}

		goto LBBEGIN
	}

	return nil
}

// Do try to execute the interface,
// specify the number of attempts, and the execution interval
func Do(retries int, trier Trier, sleeps ...Sleep) error {
	if retries <= 0 {
		return nil
	}

	total := retries

LBBEGIN:

	if err := trier.Try(); err != nil {
		retries--
		if retries == 0 {
			return ErrMaxRetries
		}

		if len(sleeps) > 0 {
			if d := sleeps[0](total - retries); d > 0 {
				time.Sleep(d)
			}
		}

		goto LBBEGIN
	}

	return nil
}
