package mysql

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type TestItem struct {
	ID   int64  `db:"id,primarykey,autoincrement"`
	Code string `db:"code,size:50"`
	Name string `db:"name,size:50"`
}

func TestDB(t *testing.T) {
	db, err := NewDB(
		SetDSN("root:123456@tcp(127.0.0.1:3306)/myapp_test?charset=utf8"),
		SetTrace(true),
	)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer db.Close()

	tableName := "test_item"
	db.AddTableWithName(TestItem{}, tableName)
	err = db.CreateTablesIfNotExists()
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer db.DropTable(TestItem{})

	err = db.Insert(&TestItem{Code: "foo", Name: "bar"})
	if err != nil {
		t.Error(err.Error())
		return
	}

	var item TestItem
	err = db.SelectOne(&item, fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName))
	if err != nil {
		t.Error(err.Error())
		return
	}

	if item.ID != 1 || item.Code != "foo" || item.Name != "bar" {
		t.Error("数据错误：", item)
		return
	}
}
