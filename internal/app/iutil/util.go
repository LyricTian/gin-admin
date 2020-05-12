package iutil

import (
	"github.com/LyricTian/gin-admin/v6/internal/app/config"
	"github.com/LyricTian/gin-admin/v6/pkg/unique"
)

var idFunc = func() string {
	return unique.NewSnowflakeID().String()
}

// InitID ...
func InitID() {
	switch config.C.UniqueID.Type {
	case "uuid":
		idFunc = func() string {
			return unique.MustUUID().String()
		}
	case "object":
		idFunc = func() string {
			return unique.NewObjectID().Hex()
		}
	default:
		idFunc = func() string {
			return unique.NewSnowflakeID().String()
		}
	}
}

// NewID Create unique id
func NewID() string {
	return idFunc()
}
