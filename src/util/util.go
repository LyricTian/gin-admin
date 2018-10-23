package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// StructToMap 将结构体转换为字典
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// StructsToMapSlice 将结构体切片转换为字典切片
func StructsToMapSlice(v interface{}) []map[string]interface{} {
	iVal := reflect.Indirect(reflect.ValueOf(v))
	if iVal.IsNil() || iVal.IsValid() || iVal.Type().Kind() != reflect.Slice {
		return make([]map[string]interface{}, 0)
	}

	l := iVal.Len()
	result := make([]map[string]interface{}, l)
	for i := 0; i < l; i++ {
		result[i] = structs.Map(iVal.Index(i).Interface())
	}

	return result
}

// Trim 去除空格
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// GetLevelCode 获取分级码
func GetLevelCode(orderLevelCodes []string) string {
	l := len(orderLevelCodes)

	if l == 0 {
		return "01"
	} else if l == 1 {
		return orderLevelCodes[0] + "01"
	}

	root := orderLevelCodes[0]

	toValue := func(i int) string {
		if i < 10 {
			return fmt.Sprintf("%s0%d", root, i)
		}
		return fmt.Sprintf("%s%d", root, i)
	}

	for i := 1; i < 100; i++ {
		code := toValue(i)
		if i <= l &&
			orderLevelCodes[i] == code {
			continue
		}
		return code
	}

	return ""
}
