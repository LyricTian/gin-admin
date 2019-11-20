package hook

import (
	"fmt"
	"os"

	"github.com/LyricTian/queue"
	"github.com/sirupsen/logrus"
)

var defaultOptions = options{
	maxQueues:  512,
	maxWorkers: 1,
	levels: []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	},
}

// ExecCloser write the logrus entry to the store and close the store
type ExecCloser interface {
	Exec(entry *logrus.Entry) error
	Close() error
}

// FilterHandle a filter handler
type FilterHandle func(*logrus.Entry) *logrus.Entry

type options struct {
	maxQueues  int
	maxWorkers int
	extra      map[string]interface{}
	filter     FilterHandle
	levels     []logrus.Level
}

// SetMaxQueues set the number of buffers
func SetMaxQueues(maxQueues int) Option {
	return func(o *options) {
		o.maxQueues = maxQueues
	}
}

// SetMaxWorkers set the number of worker threads
func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

// SetExtra set extended parameters
func SetExtra(extra map[string]interface{}) Option {
	return func(o *options) {
		o.extra = extra
	}
}

// SetFilter set the entry filter
func SetFilter(filter FilterHandle) Option {
	return func(o *options) {
		o.filter = filter
	}
}

// SetLevels set the available log level
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

// New creates a hook to be added to an instance of logger
func New(exec ExecCloser, opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
		e:    exec,
	}
}

// Hook to send logs to a mongo database
type Hook struct {
	opts options
	q    *queue.Queue
	e    ExecCloser
}

// Levels returns the available logging levels
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// Fire is called when a log event is fired
func (h *Hook) Fire(entry *logrus.Entry) error {
	entry = h.copyEntry(entry)
	h.q.Push(queue.NewJob(entry, func(v interface{}) {
		h.exec(v.(*logrus.Entry))
	}))
	return nil
}

func (h *Hook) copyEntry(e *logrus.Entry) *logrus.Entry {
	entry := logrus.NewEntry(e.Logger)
	entry.Data = make(logrus.Fields)
	entry.Time = e.Time
	entry.Level = e.Level
	entry.Message = e.Message
	for k, v := range e.Data {
		entry.Data[k] = v
	}
	return entry
}

func (h *Hook) exec(entry *logrus.Entry) {
	for k, v := range h.opts.extra {
		if _, ok := entry.Data[k]; !ok {
			entry.Data[k] = v
		}
	}

	if filter := h.opts.filter; filter != nil {
		entry = filter(entry)
	}

	err := h.e.Exec(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logrus-hook] execution error: %s", err.Error())
	}
}

// Flush waits for the log queue to be empty
func (h *Hook) Flush() {
	h.q.Terminate()
	h.e.Close()
}
