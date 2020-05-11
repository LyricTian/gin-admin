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

// StructMapToStruct 结构体映射
func StructMapToStruct(s, ts interface{}) error {
	if !structs.IsStruct(s) || !structs.IsStruct(ts) {
		return nil
	}

	ss, tss := structs.New(s), structs.New(ts)
	for _, tfield := range tss.Fields() {
		if !tfield.IsExported() {
			continue
		}

		if tfield.IsEmbedded() && tfield.Kind() == reflect.Struct {
			for _, tefield := range tfield.Fields() {
				if f, ok := ss.FieldOk(tefield.Name()); ok && !f.IsZero() {
					tefield.Set2(f.Value())
				}
			}
		} else if f, ok := ss.FieldOk(tfield.Name()); ok && !f.IsZero() {
			tfield.Set2(f.Value())
		}
	}

	return nil
}
