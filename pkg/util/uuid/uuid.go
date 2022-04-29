package uuid

import (
	"github.com/google/uuid"
)

// Define alias
type UUID = uuid.UUID

// Create uuid
func NewUUID() (UUID, error) {
	return uuid.NewRandom()
}

// Create uuid(Throw panic if something goes wrong)
func MustUUID() UUID {
	v, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return v
}

// Create uuid
func MustString() string {
	return MustUUID().String()
}
