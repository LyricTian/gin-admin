package mysqlhook

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/LyricTian/queue"
	"github.com/sirupsen/logrus"
)

var defaultOptions = options{
	out:        os.Stderr,
	maxQueues:  512,
	maxWorkers: 2,
	levels: []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	},
}

// FilterHandle a filter handler
type FilterHandle func(*logrus.Entry) *logrus.Entry

type options struct {
	maxQueues  int
	maxWorkers int
	extra      map[string]interface{}
	exec       Execer
	filter     FilterHandle
	levels     []logrus.Level
	out        io.Writer
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

// SetExec set the Execer interface
func SetExec(exec Execer) Option {
	return func(o *options) {
		o.exec = exec
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

// SetOut set error output
func SetOut(out io.Writer) Option {
	return func(o *options) {
		o.out = out
	}
}

// Option a hook parameter options
type Option func(*options)

// Default create a default mysql hook
func Default(db *sql.DB, tableName string, opts ...Option) *Hook {
	return DefaultWithExtra(db, tableName, nil, opts...)
}

// DefaultWithExtra create a default mysql hook with extra items
func DefaultWithExtra(db *sql.DB, tableName string, extraItems []*ExecExtraItem, opts ...Option) *Hook {
	var options []Option
	options = append(options, opts...)
	options = append(options, SetExec(NewExec(db, tableName, extraItems...)))
	return New(options...)
}

// New creates a hook to be added to an instance of logger
func New(opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	if opts.exec == nil {
		panic("Unknown Execer interface implementation")
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
	}
}

// Hook to send logs to a mysql database
type Hook struct {
	opts options
	q    *queue.Queue
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
	if extra := h.opts.extra; extra != nil {
		for k, v := range extra {
			if _, ok := entry.Data[k]; !ok {
				entry.Data[k] = v
			}
		}
	}

	if filter := h.opts.filter; filter != nil {
		entry = filter(entry)
	}

	err := h.opts.exec.Exec(entry)
	if err != nil && h.opts.out != nil {
		fmt.Fprintf(h.opts.out, "[Mongo-Hook] Execution error: %s", err.Error())
	}
}

// Flush waits for the log queue to be empty
func (h *Hook) Flush() {
	h.q.Terminate()
}
