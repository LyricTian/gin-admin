package utils

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var ipRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomizedIPAddr() string {
	ipRaw := make([]byte, 4)
	binary.LittleEndian.PutUint32(ipRaw, ipRand.Uint32())

	ips := make([]string, len(ipRaw))
	for i, b := range ipRaw {
		ips[i] = strconv.FormatInt(int64(b), 10)
	}
	return strings.Join(ips, ".")
}

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

func DefaultStr(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}
