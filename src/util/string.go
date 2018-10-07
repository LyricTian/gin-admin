package util

import "strconv"

// S 字符串类型转换
type S string

func (s S) String() string {
	return string(s)
}

// Int64 转换为int64
func (s S) Int64() int64 {
	i, err := strconv.ParseInt(s.String(), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Int 转换为int
func (s S) Int() int {
	return int(s.Int64())
}

// Uint 转换为uint
func (s S) Uint() uint {
	return uint(s.Uint64())
}

// Uint64 转换为uint64
func (s S) Uint64() uint64 {
	i, err := strconv.ParseUint(s.String(), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Float64 转换为float64
func (s S) Float64() float64 {
	f, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		return 0
	}
	return f
}
