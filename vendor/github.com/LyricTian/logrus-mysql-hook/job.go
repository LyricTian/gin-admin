package mysqlhook

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type job struct {
	out    io.Writer
	entry  *logrus.Entry
	exec   Execer
	extra  map[string]interface{}
	filter FilterHandle
}

func (j *job) Reset(e *logrus.Entry) {
	entry := logrus.NewEntry(e.Logger)
	entry.Data = make(logrus.Fields)
	entry.Time = e.Time
	entry.Level = e.Level
	entry.Message = e.Message
	for k, v := range e.Data {
		entry.Data[k] = v
	}
	j.entry = entry
}

func (j *job) Job() {
	if j.extra != nil {
		for k, v := range j.extra {
			if _, ok := j.entry.Data[k]; !ok {
				j.entry.Data[k] = v
			}
		}
	}

	if j.filter != nil {
		j.entry = j.filter(j.entry)
	}

	err := j.exec.Exec(j.entry)
	if err != nil && j.out != nil {
		fmt.Fprintf(j.out, "[MySQL-Hook] Execution error:%v", err)
	}
}
