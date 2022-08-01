package trace

import (
	"fmt"

	"github.com/rs/xid"
)

// New trace id
func NewTraceID() string {
	return fmt.Sprintf("trace-%s", xid.New().String())
}
