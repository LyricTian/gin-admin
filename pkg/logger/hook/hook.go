package hook

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

// Write the logrus entry to the store and close the store
type ExecCloser interface {
	Exec(entry *logrus.Entry) error
	Close() error
}

type options struct {
	maxJobs    int
	maxWorkers int
	extra      map[string]interface{}
	levels     []logrus.Level
}

// Set the number of buffers
func SetMaxJobs(maxJobs int) Option {
	return func(o *options) {
		o.maxJobs = maxJobs
	}
}

// Set the number of worker threads
func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

// Set extended parameters
func SetExtra(extra map[string]interface{}) Option {
	return func(o *options) {
		o.extra = extra
	}
}

// Set the available log level
func SetLevels(levels ...logrus.Level) Option {
	return func(o *options) {
		if len(levels) == 0 {
			return
		}
		o.levels = levels
	}
}

// Option a hook parameter options
type Option func(*options)

// Creates a hook to be added to an instance of logger
func New(exec ExecCloser, opt ...Option) *Hook {
	opts := options{
		maxJobs:    1024,
		maxWorkers: 2,
	}

	for _, o := range opt {
		o(&opts)
	}

	wg := new(sync.WaitGroup)
	wg.Add(opts.maxWorkers)

	h := &Hook{
		opts: &opts,
		q:    make(chan *logrus.Entry, opts.maxJobs),
		wg:   wg,
		e:    exec,
	}
	h.dispatch()
	return h
}

// Hook to send logs to a mongo database
type Hook struct {
	opts   *options
	q      chan *logrus.Entry
	wg     *sync.WaitGroup
	e      ExecCloser
	closed int32
}

func (h *Hook) dispatch() {
	for i := 0; i < h.opts.maxWorkers; i++ {
		go func() {
			defer func() {
				h.wg.Done()
				if r := recover(); r != nil {
					fmt.Fprintln(os.Stderr, "Hook panic:", r)
				}
			}()

			for entry := range h.q {
				err := h.e.Exec(entry)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to write entry:", err)
				}
			}
		}()
	}
}

// Returns the available logging levels
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// Fire is called when a log event is fired
func (h *Hook) Fire(entry *logrus.Entry) error {
	if atomic.LoadInt32(&h.closed) == 1 {
		return nil
	}

	if len(h.q) == h.opts.maxJobs {
		fmt.Fprintln(os.Stderr, "Too many jobs, waiting for queue to be empty, discard")
		return nil
	}

	dup := entry.Dup()
	dup.Level = entry.Level
	dup.Message = entry.Message
	for k, v := range h.opts.extra {
		if _, ok := dup.Data[k]; !ok {
			dup.Data[k] = v
		}
	}
	h.q <- dup
	return nil
}

// Waits for the log queue to be empty
func (h *Hook) Flush() {
	if atomic.LoadInt32(&h.closed) == 1 {
		return
	}
	atomic.StoreInt32(&h.closed, 1)
	close(h.q)
	h.wg.Wait()
	h.e.Close()
}
