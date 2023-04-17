package rand

import (
	"strconv"
	"testing"
)

func TestRandom(t *testing.T) {
	digits, err := Random(6, Ldigit)
	if err != nil {
		t.Error(err.Error())
		return
	} else if len(digits) != 6 {
		t.Error("invalid digit:", digits)
		return
	}

	for _, b := range digits {
		d, err := strconv.Atoi(string(b))
		if err != nil {
			t.Error(err.Error())
			return
		} else if d > 10 || d < 0 {
			t.Error("invalid digit:", d)
		}
	}
}
