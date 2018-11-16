package bll

import (
	"context"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Role 角色管理
type Role struct {
	RoleModel model.IRole `inject:"IRole"`
	MenuModel model.IMenu `inject:"IMenu"`
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
	item, err := a.RoleModel.Get(ctx, recordID, true)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrNotFound
	}

	return item, nil
}

// 过滤叶子节点
func (a *Role) filterLeafMenuIDs(ctx context.Context, menuIDs []string) ([]string, error) {
	menus, err := a.MenuModel.QuerySelect(ctx, schema.MenuSelectQueryParam{
		RecordIDs: menuIDs,
		Status:    1,
	})
	if err != nil {
		return nil, err
	}

	var leafMenuIDs []string
	for _, m := range menus {
		var exists bool
		for _, m2 := range menus {
			if strings.HasPrefix(m2.LevelCode, m.LevelCode) &&
				m2.LevelCode != m.LevelCode {
				exists = true
				break
			}
		}
		if !exists {
			leafMenuIDs = append(leafMenuIDs, m.RecordID)
		}
	}

	return leafMenuIDs, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item *schema.Role) error {
	exists, err := a.RoleModel.CheckName(ctx, item.Name)
	if err != nil {
		return err
	} else if exists {
		return errors.New("角色名称已经存在")
	}

	leafMenuIDs, err := a.filterLeafMenuIDs(ctx, item.MenuIDs)
	if err != nil {
		return err
	}
	item.MenuIDs = leafMenuIDs

	item.ID = 0
	item.RecordID = util.MustUUID()
	item.Created = time.Now().Unix()
	item.Deleted = 0
	return a.RoleModel.Create(ctx, item)
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item *schema.Role) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID, false)
	if err != nil {
		return err
	} else if oldItem == nil {
		return util.ErrNotFound
	} else if oldItem.Name != item.Name {
		exists, err := a.RoleModel.CheckName(ctx, item.Name)
		if err != nil {
			return err
		} else if exists {
			return errors.New("角色名称已经存在")
		}
	}

	leafMenuIDs, err := a.filterLeafMenuIDs(ctx, item.MenuIDs)
	if err != nil {
		return err
	}
	item.MenuIDs = leafMenuIDs

	info := util.StructToMap(item)
	delete(info, "id")
	delete(info, "record_id")
	delete(info, "creator")
	delete(info, "created")
	delete(info, "updated")
	delete(info, "deleted")

	return a.RoleModel.UpdateWithMenuIDs(ctx, recordID, info, item.MenuIDs)
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	exists, err := a.RoleModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	return a.RoleModel.Delete(ctx, recordID)
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, recordID string, status int) error {
	exists, err := a.RoleModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	info := map[string]interface{}{
		"status": status,
	}
	return a.RoleModel.Update(ctx, recordID, info)
}
