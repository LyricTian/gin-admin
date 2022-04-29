package hook

import (
	"fmt"
	"os"

	"github.com/LyricTian/queue"
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

	q := queue.NewQueue(opts.maxJobs, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: &opts,
		q:    q,
		e:    exec,
	}
}

// Hook to send logs to a mongo database
type Hook struct {
	opts *options
	q    *queue.Queue
	e    ExecCloser
}

// Returns the available logging levels
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// Fire is called when a log event is fired
func (h *Hook) Fire(entry *logrus.Entry) error {
	if h.q.GetJobCount() == h.opts.maxJobs {
		fmt.Fprintf(os.Stderr, "Too many jobs, waiting for queue to be empty, discard\n")
		return nil
	}

	dup := entry.Dup()
	dup.Level = entry.Level
	dup.Message = entry.Message
	h.q.Push(queue.NewJob(dup, func(v interface{}) {
		for k, v := range h.opts.extra {
			if _, ok := entry.Data[k]; !ok {
				entry.Data[k] = v
			}
		}
		err := h.e.Exec(entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write entry: %v\n", err)
		}
	}))
	return nil
}

// Waits for the log queue to be empty
func (h *Hook) Flush() {
	h.q.Terminate()
	h.e.Close()
}
