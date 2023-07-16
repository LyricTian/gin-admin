package util

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// The RandomizedIPAddr function generates a random IP address.
func RandomizedIPAddr() string {
	raw := make([]byte, 4)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	binary.LittleEndian.PutUint32(raw, rd.Uint32())

	ips := make([]string, len(raw))
	for i, b := range raw {
		ips[i] = strconv.FormatInt(int64(b), 10)
	}
	return strings.Join(ips, ".")
}
