package util

import (
	"fmt"
	"os"
	"time"
)

var (
	pid = os.Getpid()
)

// NewTraceID 创建追踪ID
func NewTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s",
		pid,
		time.Now().Format("2006.01.02.15.04.05.999999"))
}
