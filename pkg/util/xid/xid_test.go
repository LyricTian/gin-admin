package xid

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	t.Logf("xid: %s", strings.ToUpper(NewID()))
}
