package util

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/LyricTian/structs"
)

var (
	pid = os.Getpid()
)

// NewTraceID 创建追踪ID
func NewTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s",
		pid,
		time.Now().Format("2006.01.02.15.04.05.999999"))
}

// NewRecordID 创建记录ID
func NewRecordID() string {
	return NewObjectID().Hex()
}

// StructMapToStruct 结构体映射
func StructMapToStruct(s, ts interface{}) error {
	if !structs.IsStruct(s) || !structs.IsStruct(ts) {
		return nil
	}

	ss, tss := structs.New(s), structs.New(ts)

	var setValue = func(field *structs.Field) error {
		if sf, ok := ss.FieldOk(field.Name()); ok {
			err := field.Set2(sf.Value())
			if err != nil {
				fmt.Printf("[warning] StructMapToStruct set field [%s->%s]: %s", field.Name(), sf.Name(), err.Error())
			}
		}
		return nil
	}

	for _, field := range tss.Fields() {
		if !field.IsExported() {
			continue
		}

		if field.IsEmbedded() && field.Kind() == reflect.Struct {
			for _, field := range field.Fields() {
				setValue(field)
			}
			continue
		}
		setValue(field)
	}

	return nil
}
