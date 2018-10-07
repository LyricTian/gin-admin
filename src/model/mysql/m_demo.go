package mysql

import (
	"database/sql"
	"fmt"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/service/mysql"

	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Demo 示例model
type Demo struct {
	DB     *mysql.DB
	Common *Common
}

// Init 初始化
func (a *Demo) Init(g *inject.Graph, db *mysql.DB, c *Common) *Demo {
	a.DB = db
	a.Common = c

	g.Provide(&inject.Object{Value: model.IDemo(a), Name: "IDemo"})
	db.AddTableWithName(schema.Demo{}, a.TableName())
	return a
}

// TableName 表名
func (a *Demo) TableName() string {
	return fmt.Sprintf("%s_%s", viper.GetString("mysql_table_prefix"), "demo")
}

// Query 查询示例数据
func (a *Demo) Query() ([]*schema.Demo, error) {
	var items []*schema.Demo
	_, err := a.DB.Select(&items, fmt.Sprintf("SELECT * FROM %s", a.TableName()))
	if err != nil {
		return nil, errors.Wrap(err, "查询示例数据发生错误")
	}
	return items, nil
}

// Get 获取单条示例数据
func (a *Demo) Get(id int64) (*schema.Demo, error) {
	var item schema.Demo
	err := a.DB.SelectOne(&item, fmt.Sprintf("SELECT * FROM %s WHERE id=?", a.TableName()), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "获取单条示例数据发生错误")
	}
	return &item, nil
}

// Create 增加示例数据
func (a *Demo) Create(item *schema.Demo) error {
	err := a.DB.Insert(item)
	if err != nil {
		return errors.Wrap(err, "增加示例数据发生错误")
	}
	return nil
}

// Update 更新实例数据
func (a *Demo) Update(id int64, info map[string]interface{}) error {
	_, err := a.DB.UpdateByPK(a.TableName(), mysql.M{"id": id}, info)
	if err != nil {
		return errors.Wrap(err, "更新示例数据发生错误")
	}
	return nil
}

// Delete 删除示例数据
func (a *Demo) Delete(id int64) error {
	_, err := a.DB.DeleteByPK(a.TableName(), mysql.M{"id": id})
	if err != nil {
		return errors.Wrap(err, "删除示例数据发生错误")
	}
	return nil
}
