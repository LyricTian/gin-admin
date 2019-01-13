package util

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/fatih/structs"
)

// MustUUID 创建一个UUID，如果有错误，则抛出panic
func MustUUID() string {
	uuid, err := NewUUID()
	if err != nil {
		panic("创建UUID发生错误: " + err.Error())
	}
	return uuid
}

// NewUUID 创建UUID，参考：https://github.com/google/uuid
func NewUUID() (string, error) {
	var buf [16]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return "", err
	}

	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80

	dst := make([]byte, 36)
	hex.Encode(dst, buf[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], buf[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], buf[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], buf[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], buf[10:])

	return string(dst), nil
}

// StructToMap 将结构体转换为字典
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// StructsToMapSlice 将结构体切片转换为字典切片
func StructsToMapSlice(v interface{}) []map[string]interface{} {
	iVal := reflect.Indirect(reflect.ValueOf(v))
	if iVal.IsNil() ||
		!iVal.IsValid() ||
		iVal.Type().Kind() != reflect.Slice {
		return make([]map[string]interface{}, 0)
	}

	l := iVal.Len()
	result := make([]map[string]interface{}, l)
	for i := 0; i < l; i++ {
		result[i] = structs.Map(iVal.Index(i).Interface())
	}

	return result
}

// FillStruct 填充结构体
func FillStruct(source, target interface{}) error {
	if !structs.IsStruct(source) ||
		!structs.IsStruct(target) {
		return errors.New("invalid struct")
	} else if tv := reflect.ValueOf(target); tv.Kind() != reflect.Ptr {
		return errors.New("the target struct must be a pointer type")
	}

	sFields := structs.Fields(source)
	tFields := structs.Fields(target)

	for _, tfield := range tFields {
		for _, sfield := range sFields {
			if tfield.Name() == sfield.Name() {
				tfield.Set(sfield.Value())
				break
			}
		}
	}

	return nil
}

// FillStructs 填充结构体的切片
func FillStructs(source, target interface{}) error {
	var check = func(v reflect.Value) bool {
		if v.IsNil() ||
			!v.IsValid() ||
			v.Type().Kind() != reflect.Slice {
			return false
		}
		return true
	}

	sValue := reflect.Indirect(reflect.ValueOf(source))
	tValue := reflect.Indirect(reflect.ValueOf(target))
	if !check(sValue) || !check(tValue) ||
		sValue.Len() != tValue.Len() {
		return errors.New("invalid struct slice")
	}

	for i := 0; i < tValue.Len(); i++ {
		sv := sValue.Index(i).Interface()
		tv := tValue.Index(i)
		if tv.Kind() != reflect.Ptr {
			tv = tv.Addr()
			FillStruct(sv, tv.Interface())
		} else if tv.IsNil() {
			tv = reflect.New(tv.Type().Elem())
			FillStruct(sv, tv.Interface())
			tValue.Index(i).Set(tv)
		}
	}

	return nil
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
		if i < l &&
			orderLevelCodes[i] == code {
			continue
		}
		return code
	}

	return ""
}

// ParseLevelCodes 解析分级码（去重）
func ParseLevelCodes(levelCodes ...string) []string {
	var allCodes []string

	for _, levelCode := range levelCodes {
		codes := parseLevelCode(levelCode)

		for _, code := range codes {
			var exists bool
			for _, c := range allCodes {
				if code == c {
					exists = true
					break
				}
			}

			if !exists {
				allCodes = append(allCodes, code)
			}
		}
	}

	return allCodes
}

func parseLevelCode(levelCode string) []string {
	if len(levelCode) < 2 {
		return nil
	}
	var (
		codes []string
		root  bytes.Buffer
	)

	for i := range levelCode {
		idx := i + 1
		if idx%2 == 0 {
			root.WriteString(levelCode[idx-2 : idx])
			codes = append(codes, root.String())
		}
	}

	root.Reset()
	return codes
}

// CheckPrefix 检查是否存在前缀
func CheckPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
