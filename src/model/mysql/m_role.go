package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/service/mysql"
	"time"

	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
)

// Role 角色管理
type Role struct {
	DB     *mysql.DB
	Common *Common
}

// Init 初始化
func (a *Role) Init(g *inject.Graph, db *mysql.DB, c *Common) *Role {
	a.DB = db
	a.Common = c

	g.Provide(&inject.Object{Value: model.IRole(a), Name: "IRole"})

	db.CreateTableIfNotExists(schema.Role{}, a.TableName())
	db.CreateTableIfNotExists(schema.RoleMenu{}, a.RoleMenuTableName())

	db.CreateTableIndex(a.TableName(), "idx_record_id", true, "record_id")
	db.CreateTableIndex(a.TableName(), "idx_name", false, "name")
	db.CreateTableIndex(a.TableName(), "idx_status", false, "status")
	db.CreateTableIndex(a.TableName(), "idx_deleted", false, "deleted")
	db.CreateTableIndex(a.RoleMenuTableName(), "idx_role_id", false, "role_id")
	db.CreateTableIndex(a.RoleMenuTableName(), "idx_deleted", false, "deleted")

	return a
}

// TableName 角色表名
func (a *Role) TableName() string {
	return a.Common.TableName("role")
}

// RoleMenuTableName 角色菜单表名
func (a *Role) RoleMenuTableName() string {
	return a.Common.TableName("role_menu")
}

// QueryPage 查询分页数据
func (a *Role) QueryPage(ctx context.Context, params schema.RoleQueryParam, pageIndex, pageSize uint) (int64, []*schema.RoleQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if params.Name != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+params.Name+"%")
	}

	if params.Status != 0 {
		where = fmt.Sprintf("%s AND status = ?", where)
		args = append(args, params.Status)
	}

	count, err := a.DB.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", a.TableName(), where), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	} else if count == 0 {
		return 0, nil, nil
	}

	var items []*schema.RoleQueryResult
	fields := "id,record_id,name,memo,status"
	_, err = a.DB.Select(&items, fmt.Sprintf("SELECT %s FROM %s %s ORDER BY id DESC LIMIT %d,%d", fields, a.TableName(), where, (pageIndex-1)*pageSize, pageSize), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	}

	return count, items, nil
}

// QuerySelect 查询选择数据
func (a *Role) QuerySelect(ctx context.Context, params schema.RoleSelectQueryParam) ([]*schema.RoleSelectQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if params.Name != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+params.Name+"%")
	}

	if params.Status != 0 {
		where = fmt.Sprintf("%s AND status = ?", where)
		args = append(args, params.Status)
	}

	if len(params.RecordIDs) > 0 {
		where = fmt.Sprintf("%s AND record_id IN(?)", where)
		args = append(args, params.RecordIDs)
	}

	query := fmt.Sprintf("SELECT record_id,name FROM %s %s", a.TableName(), where)
	query, args, err := a.DB.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "查询选择数据发生错误")
	}

	var items []*schema.RoleSelectQueryResult
	_, err = a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "查询选择数据发生错误")
	}
	return items, nil
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string, includeMenuIDs bool) (*schema.Role, error) {
	var item schema.Role
	fields := "id,record_id,name,memo,status,creator,created,deleted"

	err := a.DB.SelectOne(&item, fmt.Sprintf("SELECT %s FROM %s WHERE deleted=0 AND record_id=?", fields, a.TableName()), recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "查询指定数据发生错误")
	}

	if includeMenuIDs {
		menuIDs, err := a.QueryMenuIDs(ctx, recordID)
		if err != nil {
			return nil, err
		}
		item.MenuIDs = menuIDs
	}

	return &item, nil
}

// QueryMenuIDs 查询角色菜单
func (a *Role) QueryMenuIDs(ctx context.Context, roleID string) ([]string, error) {
	query := fmt.Sprintf("SELECT menu_id FROM %s WHERE deleted=0 AND role_id=?", a.RoleMenuTableName())

	var items []*schema.RoleMenu
	_, err := a.DB.Select(&items, query, roleID)
	if err != nil {
		return nil, errors.Wrap(err, "查询角色菜单发生错误")
	}

	menuIDs := make([]string, len(items))
	for i, item := range items {
		menuIDs[i] = item.MenuID
	}

	return menuIDs, nil
}

// CheckName 检查名称
func (a *Role) CheckName(ctx context.Context, name string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted=0 AND name=?", a.TableName())
	n, err := a.DB.SelectInt(query, name)
	if err != nil {
		return false, errors.Wrap(err, "检查名称发生错误")
	}
	return n > 0, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item *schema.Role) error {
	tran, err := a.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "创建数据发生错误")
	}

	err = tran.Insert(item)
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "创建数据发生错误")
	}

	for _, menuID := range item.MenuIDs {
		roleMenuItem := &schema.RoleMenu{
			RoleID: item.RecordID,
			MenuID: menuID,
		}
		err = tran.Insert(roleMenuItem)
		if err != nil {
			tran.Rollback()
			return errors.Wrap(err, "创建数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrap(err, "创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, info map[string]interface{}) error {
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

// UpdateWithMenuIDs 更新数据
func (a *Role) UpdateWithMenuIDs(ctx context.Context, recordID string, info map[string]interface{}, menuIDs []string) error {
	tran, err := a.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "更新数据发生错误")
	}

	_, err = a.DB.UpdateByPKWithTran(tran, a.TableName(),
		map[string]interface{}{"record_id": recordID},
		info)
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "更新数据发生错误")
	}

	_, err = a.DB.UpdateByPKWithTran(tran, a.RoleMenuTableName(),
		map[string]interface{}{"role_id": recordID},
		map[string]interface{}{"deleted": time.Now().Unix()})
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "更新数据发生错误")
	}

	for _, menuID := range menuIDs {
		roleMenuItem := &schema.RoleMenu{
			RoleID: recordID,
			MenuID: menuID,
		}
		err = tran.Insert(roleMenuItem)
		if err != nil {
			tran.Rollback()
			return errors.Wrap(err, "创建数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrap(err, "更新数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	tran, err := a.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "删除数据发生错误")
	}

	_, err = a.DB.UpdateByPKWithTran(tran, a.TableName(),
		map[string]interface{}{"record_id": recordID},
		map[string]interface{}{"deleted": time.Now().Unix()})
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "删除数据发生错误")
	}

	_, err = a.DB.UpdateByPKWithTran(tran, a.RoleMenuTableName(),
		map[string]interface{}{"role_id": recordID},
		map[string]interface{}{"deleted": time.Now().Unix()})
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "删除数据发生错误")
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrap(err, "删除数据发生错误")
	}

	return nil
}
