package util

import (
	"strings"
	"testing"
)

func TestNewXID(t *testing.T) {
	t.Logf("xid: %s", strings.ToUpper(NewXID()))
}

func TestMustNewUUID(t *testing.T) {
	t.Logf("uuid: %s", strings.ToUpper(MustNewUUID()))
}
