package util

import "strconv"

func DefaultStrToInt(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func DefaultStrToBool(s string, def bool) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return b
}

func DefaultStrToFLoat(s string, def float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return f
}

func DefaultStr(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}
