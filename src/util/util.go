package util

import (
	"strings"

	"github.com/fatih/structs"
)

// StructToMap 将结构体转换为字典
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// Trim 去除空格
func Trim(s string) string {
	return strings.TrimSpace(s)
}
