package trace

import (
	"github.com/rs/xid"
)

// New trace id
func NewTraceID() string {
	return xid.New().String()
}
