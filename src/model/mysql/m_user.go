package mysql

import (
	"context"
	"fmt"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/service/mysql"
	"time"

	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
)

// User 用户管理
type User struct {
	DB     *mysql.DB
	Common *Common
}

// Init 初始化
func (a *User) Init(g *inject.Graph, db *mysql.DB, c *Common) *User {
	a.DB = db
	a.Common = c

	g.Provide(&inject.Object{Value: model.IUser(a), Name: "IUser"})

	db.CreateTableIfNotExists(schema.User{}, a.TableName())
	db.CreateTableIfNotExists(schema.UserRole{}, a.UserRoleTableName())

	db.CreateTableIndex(a.TableName(), "idx_record_id", true, "record_id")
	db.CreateTableIndex(a.TableName(), "idx_user_name", false, "user_name")
	db.CreateTableIndex(a.TableName(), "idx_real_name", false, "real_name")
	db.CreateTableIndex(a.TableName(), "idx_status", false, "status")
	db.CreateTableIndex(a.TableName(), "idx_deleted", false, "deleted")
	db.CreateTableIndex(a.UserRoleTableName(), "idx_user_id", false, "user_id")
	db.CreateTableIndex(a.UserRoleTableName(), "idx_deleted", false, "deleted")

	return a
}

// TableName 表名
func (a *User) TableName() string {
	return a.Common.TableName("user")
}

// UserRoleTableName 用户角色表名
func (a *User) UserRoleTableName() string {
	return a.Common.TableName("user_role")
}

// QueryPage 查询分页数据
func (a *User) QueryPage(ctx context.Context, params schema.UserQueryParam, pageIndex, pageSize uint) (int64, []*schema.UserQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if params.UserName != "" {
		where = fmt.Sprintf("%s AND user_name LIKE ?", where)
		args = append(args, "%"+params.UserName+"%")
	}

	if params.RealName != "" {
		where = fmt.Sprintf("%s AND real_name LIKE ?", where)
		args = append(args, "%"+params.RealName+"%")
	}

	if params.Status != 0 {
		where = fmt.Sprintf("%s AND status = ?", where)
		args = append(args, params.Status)
	}

	if params.RoleID != "" {
		where = fmt.Sprintf("%s AND record_id IN(SELECT user_id FROM %s WHERE deleted=0 AND role_id=?)", where, a.UserRoleTableName())
		args = append(args, params.RoleID)
	}

	count, err := a.DB.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", a.TableName(), where), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	} else if count == 0 {
		return 0, nil, nil
	}

	var items []*schema.UserQueryResult
	fields := "id,record_id,user_name,real_name,status,created"
	_, err = a.DB.Select(&items, fmt.Sprintf("SELECT %s FROM %s %s ORDER BY id DESC LIMIT %d,%d", fields, a.TableName(), where, (pageIndex-1)*pageSize, pageSize), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	}

	return count, items, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string) (*schema.User, error) {
	var item schema.User
	fields := "id,record_id,user_name,real_name,status,creator,created,deleted"

	err := a.DB.SelectOne(&item, fmt.Sprintf("SELECT %s FROM %s WHERE deleted=0 AND record_id=?", fields, a.TableName()), recordID)
	if err != nil {
		return nil, errors.Wrap(err, "查询指定数据发生错误")
	}

	roleIDs, err := a.QueryRoleIDs(ctx, recordID)
	if err != nil {
		return nil, err
	}
	item.RoleIDs = roleIDs

	return &item, nil
}

// QueryRoleIDs 查询用户角色
func (a *User) QueryRoleIDs(ctx context.Context, userID string) ([]string, error) {
	query := fmt.Sprintf("SELECT role_id FROM %s WHERE deleted=0 AND user_id=?", a.UserRoleTableName())

	var items []*schema.UserRole
	_, err := a.DB.Select(&items, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "查询用户角色发生错误")
	}

	roleIDs := make([]string, len(items))
	for i, item := range items {
		roleIDs[i] = item.RoleID
	}

	return roleIDs, nil
}

// CheckUserName 检查用户名
func (a *User) CheckUserName(ctx context.Context, userName string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted=0 AND user_name=?", a.TableName())
	n, err := a.DB.SelectInt(query, userName)
	if err != nil {
		return false, errors.Wrap(err, "检查用户名发生错误")
	}
	return n > 0, nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item *schema.User) error {
	tran, err := a.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "创建数据发生错误")
	}

	err = tran.Insert(item)
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "创建数据发生错误")
	}

	for _, roleID := range item.RoleIDs {
		userRoleItem := &schema.UserRole{
			UserID: item.RecordID,
			RoleID: roleID,
		}
		err = tran.Insert(userRoleItem)
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
func (a *User) Update(ctx context.Context, recordID string, info map[string]interface{}) error {
	_, err := a.DB.UpdateByPK(a.TableName(),
		map[string]interface{}{"record_id": recordID},
		info)
	if err != nil {
		return errors.Wrap(err, "更新数据发生错误")
	}
	return nil
}

// UpdateWithRoleIDs 更新数据
func (a *User) UpdateWithRoleIDs(ctx context.Context, recordID string, info map[string]interface{}, roleIDs []string) error {
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

	_, err = a.DB.UpdateByPKWithTran(tran, a.UserRoleTableName(),
		map[string]interface{}{"user_id": recordID},
		map[string]interface{}{"deleted": time.Now().Unix()})
	if err != nil {
		tran.Rollback()
		return errors.Wrap(err, "更新数据发生错误")
	}

	for _, roleID := range roleIDs {
		userRoleItem := &schema.UserRole{
			UserID: recordID,
			RoleID: roleID,
		}
		err = tran.Insert(userRoleItem)
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
func (a *User) Delete(ctx context.Context, recordID string) error {
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

	_, err = a.DB.UpdateByPKWithTran(tran, a.UserRoleTableName(),
		map[string]interface{}{"user_id": recordID},
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
