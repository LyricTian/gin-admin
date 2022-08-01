package xid

import "github.com/rs/xid"

// Create global unique id for use the Object ID algorithm
func NewID() string {
	return xid.New().String()
}
