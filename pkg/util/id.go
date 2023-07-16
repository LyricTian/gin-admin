package util

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

// The function "NewXID" generates a new unique identifier (XID) and returns it as a string.
func NewXID() string {
	return xid.New().String()
}

// The function generates a new UUID and panics if there is an error.
func MustNewUUID() string {
	v, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return v.String()
}
