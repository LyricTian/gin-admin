package yaml

import (
	"gopkg.in/yaml.v2"
)

// 定义YAML操作
var (
	Marshal    = yaml.Marshal
	Unmarshal  = yaml.Unmarshal
	NewDecoder = yaml.NewDecoder
	NewEncoder = yaml.NewEncoder
)
