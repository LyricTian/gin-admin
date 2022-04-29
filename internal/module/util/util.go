package util

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v9/pkg/util/json"
	"github.com/go-playground/validator/v10"
)

// user authorization object
type Subject struct {
	UserID string `json:"u,omitempty"`
	Role   string `json:"r,omitempty"`
}

func GenerateSubject(sub Subject) string {
	return json.MarshalToString(sub)
}

func ParseSubject(str string) Subject {
	var sub Subject
	json.Unmarshal([]byte(str), &sub)
	return sub
}

var (
	ipRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RandomizedIPAddr() string {
	ipRaw := make([]byte, 4)
	binary.LittleEndian.PutUint32(ipRaw, ipRand.Uint32())

	ips := make([]string, len(ipRaw))
	for i, b := range ipRaw {
		ips[i] = strconv.FormatInt(int64(b), 10)
	}

	return strings.Join(ips, ".")
}

var (
	defaultValidate = validator.New()
)

func ValidateEmail(email string) bool {
	return defaultValidate.Var(email, "email") == nil
}
