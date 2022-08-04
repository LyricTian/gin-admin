package generate

import "testing"

func TestToLowerUnderlined(t *testing.T) {
	name := ToLowerUnderlined("HelloWorld")
	if name != "hello_world" {
		t.Errorf("ToLowerUnderlined error, expect: %s, actual: %s", "hello_world", name)
	}
}
