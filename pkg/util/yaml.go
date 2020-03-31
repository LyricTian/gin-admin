package util

import (
	"gopkg.in/yaml.v2"
)

// 定义YAML操作
var (
	YAMLMarshal    = yaml.Marshal
	YAMLUnmarshal  = yaml.Unmarshal
	YAMLNewDecoder = yaml.NewDecoder
	YAMLNewEncoder = yaml.NewEncoder
)
