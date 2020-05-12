package trace

import (
	"fmt"
	"os"
	"time"
)

// IDFunc Get trace id
type IDFunc func() string

var (
	idFunc IDFunc
	pid    = os.Getpid()
)

func init() {
	idFunc = func() string {
		return fmt.Sprintf("trace-id-%d-%s",
			pid,
			time.Now().Format("2006.01.02.15.04.05.999999"))
	}
}

// SetIDFunc Set trace id func
func SetIDFunc(traceIDFunc IDFunc) {
	idFunc = traceIDFunc
}

// NewID Create trace id
func NewID() string {
	return idFunc()
}
