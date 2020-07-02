package util

import (
	"github.com/jinzhu/copier"
)

// StructMapToStruct 结构体映射
func StructMapToStruct(s, ts interface{}) error {
	return copier.Copy(ts, s)
}
