package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/mysql"
	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
)

// Demo 示例程序
type Demo struct {
	DB     *mysql.DB
	Common *Common
}

// Init 初始化
func (a *Demo) Init(g *inject.Graph, db *mysql.DB, c *Common) *Demo {
	a.DB = db
	a.Common = c

	g.Provide(&inject.Object{Value: model.IDemo(a), Name: "IDemo"})

	db.CreateTableIfNotExists(schema.Demo{}, a.TableName())

	db.CreateTableIndex(a.TableName(), "idx_record_id", true, "record_id")
	db.CreateTableIndex(a.TableName(), "idx_code", false, "code")
	db.CreateTableIndex(a.TableName(), "idx_name", false, "name")
	db.CreateTableIndex(a.TableName(), "idx_deleted", false, "deleted")

	return a
}

// TableName 表名
func (a *Demo) TableName() string {
	return a.Common.TableName("demo")
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx context.Context, params schema.DemoQueryParam, pageIndex, pageSize uint) (int64, []*schema.DemoQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if params.Code != "" {
		where = fmt.Sprintf("%s AND code LIKE ?", where)
		args = append(args, "%"+params.Code+"%")
	}

	if params.Name != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+params.Name+"%")
	}

	count, err := a.DB.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", a.TableName(), where), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	} else if count == 0 {
		return 0, nil, nil
	}

	var items []*schema.DemoQueryResult
	fields := "id,record_id,code,name"
	_, err = a.DB.Select(&items, fmt.Sprintf("SELECT %s FROM %s %s ORDER BY id DESC LIMIT %d,%d", fields, a.TableName(), where, (pageIndex-1)*pageSize, pageSize), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	}

	return count, items, nil
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string) (*schema.Demo, error) {
	var item schema.Demo
	fields := "id,record_id,code,name,creator,created,deleted"

	err := a.DB.SelectOne(&item, fmt.Sprintf("SELECT %s FROM %s WHERE deleted=0 AND record_id=?", fields, a.TableName()), recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "查询指定数据发生错误")
	}
	return &item, nil
}

// Check 检查数据是否存在
func (a *Demo) Check(ctx context.Context, recordID string) (bool, error) {
	n, err := a.DB.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted=0 AND record_id=?", a.TableName()), recordID)
	if err != nil {
		return false, errors.Wrap(err, "检查数据是否存在发生错误")
	}

	return n > 0, nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item *schema.Demo) error {
	err := a.DB.Insert(item)
	if err != nil {
		return errors.Wrap(err, "创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, info map[string]interface{}) error {
	if _, ok := info["updated"]; !ok {
		info["updated"] = time.Now().Unix()
	}

	_, err := a.DB.UpdateByPK(a.TableName(),
		map[string]interface{}{"record_id": recordID},
		info)
	if err != nil {
		return errors.Wrap(err, "更新数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	_, err := a.DB.UpdateByPK(a.TableName(),
		map[string]interface{}{"record_id": recordID},
		map[string]interface{}{"deleted": time.Now().Unix()})
	if err != nil {
		return errors.Wrap(err, "删除数据发生错误")
	}
	return nil
}
