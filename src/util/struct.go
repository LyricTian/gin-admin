package util

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

// 定义错误
var (
	ErrInvalidStruct      = errors.New("invalid struct")
	ErrInvalidStructSlice = errors.New("invalid struct slice")
)

// StructToMap 将结构体转换为字典
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// StructsToMapSlice 将结构体切片转换为字典切片
func StructsToMapSlice(v interface{}) []map[string]interface{} {
	iVal := reflect.Indirect(reflect.ValueOf(v))
	if iVal.IsNil() ||
		iVal.Kind() != reflect.Slice {
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
func FillStruct(src, dst interface{}) error {
	if !structs.IsStruct(src) ||
		!structs.IsStruct(dst) {
		return ErrInvalidStruct
	} else if tv := reflect.ValueOf(dst); tv.Kind() != reflect.Ptr || tv.IsNil() {
		return ErrInvalidStruct
	}

	sFields := structs.Fields(src)
	tFields := structs.Fields(dst)

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

// FillStructs 填充结构体切片
func FillStructs(src, dst interface{}) error {
	sv := reflect.Indirect(reflect.ValueOf(src))
	if sv.IsNil() || sv.Kind() != reflect.Slice {
		return ErrInvalidStructSlice
	}

	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr || dv.IsNil() {
		return ErrInvalidStructSlice
	}

	sl := sv.Len()
	dv = dv.Elem()
	dt := dv.Type()
	det := dt.Elem()
	isPtr := det.Kind() == reflect.Ptr
	ndv := reflect.MakeSlice(dt, sl, sl)

	for i := 0; i < sl; i++ {
		sele := sv.Index(i).Interface()
		if isPtr {
			nv := reflect.New(det.Elem())
			FillStruct(sele, nv.Interface())
			ndv.Index(i).Set(nv)
			continue
		}
		nv := ndv.Index(i).Addr()
		FillStruct(sele, nv.Interface())
	}

	dv.Set(ndv)
	return nil
}
