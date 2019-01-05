package mysqlhook

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Execer write the logrus entry to the database
type Execer interface {
	Exec(entry *logrus.Entry) error
}

// NewExecExtraItem create extra item instance
func NewExecExtraItem(field, dbType string) *ExecExtraItem {
	return &ExecExtraItem{
		Field:  field,
		DBType: dbType,
	}
}

// ExecExtraItem extra item
type ExecExtraItem struct {
	Field  string // field name
	DBType string // mysql data type
}

// NewExec create an exec instance
func NewExec(db *sql.DB, tableName string, extraItems ...*ExecExtraItem) Execer {
	var sourceItems []*ExecExtraItem
	sourceItems = append(sourceItems, NewExecExtraItem("id", "bigint not null primary key auto_increment"))
	sourceItems = append(sourceItems, NewExecExtraItem("level", "int"))
	sourceItems = append(sourceItems, NewExecExtraItem("message", "varchar(1024)"))
	if len(extraItems) > 0 {
		sourceItems = append(sourceItems, extraItems...)
	}
	sourceItems = append(sourceItems, NewExecExtraItem("data", "text"))
	sourceItems = append(sourceItems, NewExecExtraItem("time", "DATETIME"))

	var fields []string
	for _, item := range sourceItems {
		fields = append(fields, fmt.Sprintf("%s %s", item.Field, item.DBType))
	}

	query := fmt.Sprintf("create table if not exists `%s` (%s)  engine=MyISAM charset=UTF8;", tableName, strings.Join(fields, ","))
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}

	return &defaultExec{db, tableName, sourceItems}
}

type defaultExec struct {
	db         *sql.DB
	tableName  string
	extraItems []*ExecExtraItem
}

func (e *defaultExec) Exec(entry *logrus.Entry) error {
	var (
		fields       []string
		placeholders []string
		values       []interface{}
	)

	for _, item := range e.extraItems {
		if item.Field == "id" {
			continue
		}

		fields = append(fields, item.Field)
		placeholders = append(placeholders, "?")

		switch item.Field {
		case "level":
			values = append(values, entry.Level)
		case "message":
			values = append(values, entry.Message)
		case "data":
			var data string
			if len(entry.Data) > 0 {
				jsonData, _ := json.Marshal(entry.Data)
				data = string(jsonData)
			}
			values = append(values, data)
		case "time":
			values = append(values, entry.Time)
		default:
			var value interface{}
			if v, ok := entry.Data[item.Field]; ok {
				value = v
				delete(entry.Data, item.Field)
			}
			values = append(values, value)
		}
	}

	query := fmt.Sprintf("insert into `%s` (%s) values (%s);", e.tableName, strings.Join(fields, ","), strings.Join(placeholders, ","))
	_, err := e.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}
