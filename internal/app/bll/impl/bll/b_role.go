package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
)

var _ bll.IRole = (*Role)(nil)

// RoleSet 注入Role
var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"), wire.Bind(new(bll.IRole), new(*Role)))

// Role 角色管理
type Role struct {
	Enforcer      *casbin.SyncedEnforcer
	TransModel    model.ITrans
	RoleModel     model.IRole
	RoleMenuModel model.IRoleMenu
	UserModel     model.IUser
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return a.RoleModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	item, err := a.RoleModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	roleMenus, err := a.QueryRoleMenus(ctx, recordID)
	if err != nil {
		return nil, err
	}
	item.RoleMenus = roleMenus

	return item, nil
}

// QueryRoleMenus 查询角色菜单列表
func (a *Role) QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error) {
	result, err := a.RoleMenuModel.Query(ctx, schema.RoleMenuQueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.RecordIDResult, error) {
	err := a.checkName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.NewRecordID()
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		for _, rmItem := range item.RoleMenus {
			rmItem.RecordID = util.NewRecordID()
			rmItem.RoleID = item.RecordID
			err := a.RoleMenuModel.Create(ctx, *rmItem)
			if err != nil {
				return err
			}
		}
		return a.RoleModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewRecordIDResult(item.RecordID), nil
}

func (a *Role) checkName(ctx context.Context, item schema.Role) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("角色名称已经存在")
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item schema.Role) error {
	oldItem, err := a.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		err := a.checkName(ctx, item)
		if err != nil {
			return err
		}
	}

	item.RecordID = oldItem.RecordID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		addRoleMenus, delRoleMenus := a.compareRoleMenus(ctx, oldItem.RoleMenus, item.RoleMenus)
		for _, rmitem := range addRoleMenus {
			rmitem.RecordID = util.NewRecordID()
			rmitem.RoleID = recordID
			err := a.RoleMenuModel.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delRoleMenus {
			err := a.RoleMenuModel.Delete(ctx, rmitem.RecordID)
			if err != nil {
				return err
			}
		}

		return a.RoleModel.Update(ctx, recordID, item)
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

func (a *Role) compareRoleMenus(ctx context.Context, oldRoleMenus, newRoleMenus schema.RoleMenus) (addList, delList schema.RoleMenus) {
	mOldRoleMenus := oldRoleMenus.ToMap()
	mNewRoleMenus := newRoleMenus.ToMap()

	for k, item := range mNewRoleMenus {
		if _, ok := mOldRoleMenus[k]; ok {
			delete(mOldRoleMenus, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldRoleMenus {
		delList = append(delList, item)
	}
	return
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		RoleIDs:         []string{recordID},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return errors.New400Response("该角色已被赋予用户，不允许删除")
	}

	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		err := a.RoleMenuModel.DeleteByRoleID(ctx, recordID)
		if err != nil {
			return err
		}

		return a.RoleModel.Delete(ctx, recordID)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, recordID string, status int) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.RoleModel.UpdateStatus(ctx, recordID, status)
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}
