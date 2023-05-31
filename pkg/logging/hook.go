package logging

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type HookExecuter interface {
	Exec(extra map[string]string, b []byte) error
	Close() error
}

type hookOptions struct {
	maxJobs    int
	maxWorkers int
	extra      map[string]string
}

// Set the number of buffers
func SetHookMaxJobs(maxJobs int) HookOption {
	return func(o *hookOptions) {
		o.maxJobs = maxJobs
	}
}

// Set the number of worker threads
func SetHookMaxWorkers(maxWorkers int) HookOption {
	return func(o *hookOptions) {
		o.maxWorkers = maxWorkers
	}
}

// Set extended parameters
func SetHookExtra(extra map[string]string) HookOption {
	return func(o *hookOptions) {
		o.extra = extra
	}
}

// HookOption a hook parameter options
type HookOption func(*hookOptions)

// Creates a hook to be added to an instance of logger
func NewHook(exec HookExecuter, opt ...HookOption) *Hook {
	opts := &hookOptions{
		maxJobs:    1024,
		maxWorkers: 2,
	}

	for _, o := range opt {
		o(opts)
	}

	wg := new(sync.WaitGroup)
	wg.Add(opts.maxWorkers)

	h := &Hook{
		opts: opts,
		q:    make(chan []byte, opts.maxJobs),
		wg:   wg,
		e:    exec,
	}
	h.dispatch()
	return h
}

// Hook to send logs to a mongo database
type Hook struct {
	opts   *hookOptions
	q      chan []byte
	wg     *sync.WaitGroup
	e      HookExecuter
	closed int32
}

func (h *Hook) dispatch() {
	for i := 0; i < h.opts.maxWorkers; i++ {
		go func() {
			defer func() {
				h.wg.Done()
				if r := recover(); r != nil {
					fmt.Println("Recovered from panic in logger hook:", r)
				}
			}()

			for data := range h.q {
				err := h.e.Exec(h.opts.extra, data)
				if err != nil {
					fmt.Println("Failed to write entry:", err.Error())
				}
			}
		}()
	}
}

func (h *Hook) Write(p []byte) (int, error) {
	if atomic.LoadInt32(&h.closed) == 1 {
		return len(p), nil
	}
	if len(h.q) == h.opts.maxJobs {
		fmt.Println("Too many jobs, waiting for queue to be empty, discard")
		return len(p), nil
	}

	data := make([]byte, len(p))
	copy(data, p)
	h.q <- data

	return len(p), nil
}

// Waits for the log queue to be empty
func (h *Hook) Flush() {
	if atomic.LoadInt32(&h.closed) == 1 {
		return
	}
	atomic.StoreInt32(&h.closed, 1)
	close(h.q)
	h.wg.Wait()
	err := h.e.Close()
	if err != nil {
		fmt.Println("Failed to close logger hook:", err.Error())
	}
}
