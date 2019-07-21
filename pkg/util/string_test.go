package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS(t *testing.T) {
	intS := S("1010")
	i, err := intS.Int64()
	assert.Nil(t, err)
	assert.EqualValues(t, 1010, i)

	floatS := S("10.1")
	f, err := floatS.Float64()
	assert.Nil(t, err)
	assert.EqualValues(t, 10.1, f)

	di := floatS.DefaultInt64(5)
	assert.EqualValues(t, 5, di)

	boolS := S("true")
	b, err := boolS.Bool()
	assert.Nil(t, err)
	assert.EqualValues(t, true, b)
}
