package snowflake

import "testing"

func TestMustID(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(i+1, ": ", MustID())
	}
}
