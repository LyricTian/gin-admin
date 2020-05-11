package iutil

import (
	"github.com/LyricTian/gin-admin/v6/internal/app/config"
	"github.com/LyricTian/gin-admin/v6/pkg/unique"
	"github.com/bwmarrin/snowflake"
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
		snowflake.Epoch = config.C.UniqueID.Snowflake.Epoch
		node, err := snowflake.NewNode(config.C.UniqueID.Snowflake.Node)
		if err != nil {
			panic(err)
		}
		idFunc = func() string {
			return node.Generate().String()
		}
	}
}

// NewID Create unique id
func NewID() string {
	return idFunc()
}
