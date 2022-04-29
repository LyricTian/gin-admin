package uuid

import (
	"strings"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	t.Log(strings.Replace(strings.ToUpper(MustString()), "-", "", -1))
}
