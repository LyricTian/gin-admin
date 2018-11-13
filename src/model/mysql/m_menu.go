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

// Menu 菜单管理
type Menu struct {
	DB     *mysql.DB
	Common *Common
}

// Init 初始化
func (a *Menu) Init(g *inject.Graph, db *mysql.DB, c *Common) *Menu {
	a.DB = db
	a.Common = c

	g.Provide(&inject.Object{Value: model.IMenu(a), Name: "IMenu"})

	db.CreateTableIfNotExists(schema.Menu{}, a.TableName())
	db.CreateTableIndex(a.TableName(), "idx_record_id", true, "record_id")
	db.CreateTableIndex(a.TableName(), "idx_code", false, "code")
	db.CreateTableIndex(a.TableName(), "idx_name", false, "name")
	db.CreateTableIndex(a.TableName(), "idx_type", false, "type")
	db.CreateTableIndex(a.TableName(), "idx_parent_id", false, "parent_id")
	db.CreateTableIndex(a.TableName(), "idx_status", false, "status")
	db.CreateTableIndex(a.TableName(), "idx_deleted", false, "deleted")

	return a
}

// TableName 表名
func (a *Menu) TableName() string {
	return a.Common.TableName("menu")
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx context.Context, params schema.MenuQueryParam, pageIndex, pageSize uint) (int64, []*schema.MenuQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if v := params.Name; v != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+v+"%")
	}
	if v := params.ParentID; v != "" {
		where = fmt.Sprintf("%s AND parent_id=?", where)
		args = append(args, v)
	}
	if v := params.Status; v > 0 {
		where = fmt.Sprintf("%s AND status=?", where)
		args = append(args, v)
	}
	if v := params.Type; v > 0 {
		where = fmt.Sprintf("%s AND type=?", where)
		args = append(args, v)
	}

	count, err := a.DB.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", a.TableName(), where), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	} else if count == 0 {
		return 0, nil, nil
	}

	var items []*schema.MenuQueryResult
	fields := "id,record_id,code,name,icon,path,type,sequence,status"
	_, err = a.DB.Select(&items, fmt.Sprintf("SELECT %s FROM %s %s ORDER BY type,sequence,id LIMIT %d,%d", fields, a.TableName(), where, (pageIndex-1)*pageSize, pageSize), args...)
	if err != nil {
		return 0, nil, errors.Wrap(err, "查询分页数据发生错误")
	}

	return count, items, nil
}

// QuerySelect 查询选择数据
func (a *Menu) QuerySelect(ctx context.Context, params schema.MenuSelectQueryParam) ([]*schema.MenuSelectQueryResult, error) {
	var (
		where = "WHERE deleted=0"
		args  []interface{}
	)

	if v := params.Name; v != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		where = fmt.Sprintf("%s AND status=?", where)
		args = append(args, v)
	}
	if v := params.SystemCode; v != "" {
		menu, err := a.GetByCodeAndType(ctx, v, 10)
		if err != nil {
			return nil, err
		} else if menu != nil {
			where = fmt.Sprintf("%s AND level_code!=? AND level_code LIKE ?", where)
			args = append(args, menu.LevelCode, menu.LevelCode+"%")
		}
	}

	if v := params.UserID; v != "" {
		where = fmt.Sprintf("%s AND record_id IN(SELECT menu_id FROM %s WHERE deleted=0 AND role_id IN(SELECT role_id FROM %s WHERE deleted=0 AND user_id=?))",
			where,
			a.Common.Role.RoleMenuTableName(),
			a.Common.User.UserRoleTableName(),
		)
		args = append(args, v)
	}

	if v := params.RecordIDs; len(v) > 0 {
		where = fmt.Sprintf("%s AND record_id IN(?)", where)
		args = append(args, v)
	}

	var items []*schema.MenuSelectQueryResult

	fields := "record_id,code,name,level_code,parent_id,type,icon,path"
	query, args, err := a.DB.In(fmt.Sprintf("SELECT %s FROM %s %s ORDER BY sequence,id", fields, a.TableName(), where), args...)
	if err != nil {
		return nil, errors.Wrap(err, "查询选择数据发生错误")
	}

	_, err = a.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "查询选择数据发生错误")
	}

	return items, nil
}

func (a *Menu) getAllFields() string {
	fields := "id,record_id,code,name,sequence,icon,path,level_code,parent_id,status,creator,created,deleted"
	return fields
}

// GetByCodeAndType 根据编号和类型查询指定数据
func (a *Menu) GetByCodeAndType(ctx context.Context, code string, typ int) (*schema.Menu, error) {
	var item schema.Menu

	fields := a.getAllFields()
	err := a.DB.SelectOne(&item, fmt.Sprintf("SELECT %s FROM %s WHERE deleted=0 AND code=? AND type=?", fields, a.TableName()), code, typ)
	if err != nil {
		return nil, errors.Wrap(err, "根据编号和类型查询指定数据发生错误")
	}
	return &item, nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string) (*schema.Menu, error) {
	var item schema.Menu

	fields := a.getAllFields()
	err := a.DB.SelectOne(&item, fmt.Sprintf("SELECT %s FROM %s WHERE deleted=0 AND record_id=?", fields, a.TableName()), recordID)
	if err != nil {
		return nil, errors.Wrap(err, "查询指定数据发生错误")
	}
	return &item, nil
}

// CheckCode 检查编号是否存在
func (a *Menu) CheckCode(ctx context.Context, code string, parentID string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted=0 AND code=? AND parent_id=?", a.TableName())

	n, err := a.DB.SelectInt(query, code, parentID)
	if err != nil {
		return false, errors.Wrap(err, "检查编号是否存在发生错误")
	}
	return n > 0, nil
}

// QueryLevelCodesByParentID 根据父级查询分级码
func (a *Menu) QueryLevelCodesByParentID(parentID string) ([]string, error) {
	query := fmt.Sprintf("SELECT level_code FROM %s WHERE deleted=0 AND (parent_id=? OR record_id=?) ORDER BY level_code", a.TableName())

	var items []*schema.Menu
	_, err := a.DB.Select(&items, query, parentID, parentID)
	if err != nil {
		return nil, errors.Wrap(err, "根据父级查询分级码发生错误")
	}

	levelCodes := make([]string, len(items))
	for i, item := range items {
		levelCodes[i] = item.LevelCode
	}

	return levelCodes, nil
}

// CheckChild 检查子级是否存在
func (a *Menu) CheckChild(ctx context.Context, parentID string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted=0 AND parent_id=?", a.TableName())

	n, err := a.DB.SelectInt(query, parentID)
	if err != nil {
		return false, errors.Wrap(err, "检查子级是否存在发生错误")
	}
	return n > 0, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item *schema.Menu) error {
	err := a.DB.Insert(item)
	if err != nil {
		return errors.Wrap(err, "创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, info map[string]interface{}) error {
	_, err := a.DB.UpdateByPK(a.TableName(),
		map[string]interface{}{"record_id": recordID},
		info)
	if err != nil {
		return errors.Wrap(err, "更新数据发生错误")
	}
	return nil
}

// UpdateWithLevelCode 更新数据
func (a *Menu) UpdateWithLevelCode(ctx context.Context, recordID string, info map[string]interface{}, oldLevelCode, newLevelCode string) error {
	tran, err := a.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "更新数据发生错误")
	}

	_, err = a.DB.UpdateByPKWithTran(tran, a.TableName(), map[string]interface{}{"record_id": recordID}, info)
	if err != nil {
		tran.Rollback()
		return errors.Wrapf(err, "更新数据发生错误")
	}

	query := fmt.Sprintf("UPDATE %s SET level_code=concat('%s',substr(level_code,%d)) WHERE deleted=0 AND level_code LIKE '%s%%' ORDER BY level_code", a.TableName(), newLevelCode, len(oldLevelCode)+1, oldLevelCode)
	_, err = tran.Exec(query)
	if err != nil {
		tran.Rollback()
		return errors.Wrapf(err, "更新数据发生错误")
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "更新数据提交事物发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	_, err := a.DB.UpdateByPK(a.TableName(),
		map[string]interface{}{"record_id": recordID},
		map[string]interface{}{"deleted": time.Now().Unix()})
	if err != nil {
		return errors.Wrap(err, "删除数据发生错误")
	}
	return nil
}
