package idx

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

// Create global unique id for use the Object ID algorithm
func NewXID() string {
	return xid.New().String()
}

// Create global unique id for use the UUID algorithm
func MustNewUUID() string {
	v, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return v.String()
}
