package unique

import (
	"github.com/bwmarrin/snowflake"
)

// SnowflakeNode Define snowflake node (default 1)
var SnowflakeNode *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	SnowflakeNode = node
}

// SnowflakeID Define alias
type SnowflakeID = snowflake.ID

// NewSnowflakeID Create snowflake id
func NewSnowflakeID() SnowflakeID {
	return SnowflakeNode.Generate()
}
