package bll

import (
	"context"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"time"

	"github.com/pkg/errors"
)

// Role 角色管理
type Role struct {
	RoleModel model.IRole `inject:"IRole"`
}

// QueryPage 查询分页数据
func (a *Role) QueryPage(ctx context.Context, params schema.RoleQueryParam, pageIndex, pageSize uint) (int64, []*schema.RoleQueryResult, error) {
	return a.RoleModel.QueryPage(ctx, params, pageIndex, pageSize)
}

// QuerySelect 查询选择数据
func (a *Role) QuerySelect(ctx context.Context, params schema.RoleSelectQueryParam) ([]*schema.RoleSelectQueryResult, error) {
	return a.RoleModel.QuerySelect(ctx, params)
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string) (*schema.Role, error) {
	return a.RoleModel.Get(ctx, recordID)
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item *schema.Role) error {
	exists, err := a.RoleModel.CheckName(ctx, item.Name)
	if err != nil {
		return err
	} else if exists {
		return errors.New("角色名称已经存在")
	}

	item.ID = 0
	item.RecordID = util.UUIDString()
	item.Created = time.Now().Unix()
	item.Deleted = 0
	return a.RoleModel.Create(ctx, item)
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item *schema.Role) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("无效的数据")
	} else if oldItem.Name != item.Name {
		exists, err := a.RoleModel.CheckName(ctx, item.Name)
		if err != nil {
			return err
		} else if exists {
			return errors.New("角色名称已经存在")
		}
	}

	info := util.StructToMap(item)
	delete(info, "id")
	delete(info, "record_id")
	delete(info, "creator")
	delete(info, "created")
	delete(info, "deleted")

	return a.RoleModel.UpdateWithMenuIDs(ctx, recordID, info, item.MenuIDs)
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	return a.RoleModel.Delete(ctx, recordID)
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, recordID string, status int) error {
	info := map[string]interface{}{
		"status": status,
	}
	return a.RoleModel.Update(ctx, recordID, info)
}
