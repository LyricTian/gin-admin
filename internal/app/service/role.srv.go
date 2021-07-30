package service

import (
	"context"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
	"github.com/LyricTian/gin-admin/v8/pkg/util/snowflake"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleSrv), "*"))

// RoleSrv 角色管理
type RoleSrv struct {
	Enforcer               *casbin.SyncedEnforcer
	TransRepo              *dao.TransRepo
	RoleRepo               *dao.RoleRepo
	RoleMenuRepo           *dao.RoleMenuRepo
	UserRepo               *dao.UserRepo
	MenuActionResourceRepo *dao.MenuActionResourceRepo
}

// Query 查询数据
func (a *RoleSrv) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return a.RoleRepo.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *RoleSrv) Get(ctx context.Context, id uint64, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	item, err := a.RoleRepo.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	roleMenus, err := a.QueryRoleMenus(ctx, id)
	if err != nil {
		return nil, err
	}
	item.RoleMenus = roleMenus

	return item, nil
}

// QueryRoleMenus 查询角色菜单列表
func (a *RoleSrv) QueryRoleMenus(ctx context.Context, roleID uint64) (schema.RoleMenus, error) {
	result, err := a.RoleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Create 创建数据
func (a *RoleSrv) Create(ctx context.Context, item schema.Role) (*schema.IDResult, error) {
	err := a.checkName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = snowflake.MustID()
	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		for _, rmItem := range item.RoleMenus {
			rmItem.ID = snowflake.MustID()
			rmItem.RoleID = item.ID
			err := a.RoleMenuRepo.Create(ctx, *rmItem)
			if err != nil {
				return err
			}
		}
		return a.RoleRepo.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	// Sync update casbin for role
	resources, err := a.MenuActionResourceRepo.Query(ctx, schema.MenuActionResourceQueryParam{
		MenuIDs: item.RoleMenus.ToMenuIDs(),
	})
	if err != nil {
		return nil, err
	}
	for _, ritem := range resources.Data.ToMap() {
		a.Enforcer.AddPermissionForUser(strconv.FormatUint(item.ID, 10), ritem.Path, ritem.Method)
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *RoleSrv) checkName(ctx context.Context, item schema.Role) error {
	result, err := a.RoleRepo.Query(ctx, schema.RoleQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("名称不允许重复")
	}
	return nil
}

// Update 更新数据
func (a *RoleSrv) Update(ctx context.Context, id uint64, item schema.Role) error {
	oldItem, err := a.Get(ctx, id)
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

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		addRoleMenus, delRoleMenus := a.compareRoleMenus(ctx, oldItem.RoleMenus, item.RoleMenus)
		for _, rmitem := range addRoleMenus {
			rmitem.ID = snowflake.MustID()
			rmitem.RoleID = id
			err := a.RoleMenuRepo.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delRoleMenus {
			err := a.RoleMenuRepo.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.RoleRepo.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	// Sync update casbin for role
	roleMenus, err := a.RoleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{
		RoleID: id,
	})
	if err != nil {
		return err
	}

	resources, err := a.MenuActionResourceRepo.Query(ctx, schema.MenuActionResourceQueryParam{
		MenuIDs: roleMenus.Data.ToMenuIDs(),
	})
	if err != nil {
		return err
	}

	a.Enforcer.DeleteRole(strconv.FormatUint(item.ID, 10))
	for _, ritem := range resources.Data.ToMap() {
		a.Enforcer.AddPermissionForUser(strconv.FormatUint(item.ID, 10), ritem.Path, ritem.Method)
	}

	return nil
}

func (a *RoleSrv) compareRoleMenus(ctx context.Context, oldRoleMenus, newRoleMenus schema.RoleMenus) (addList, delList schema.RoleMenus) {
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
func (a *RoleSrv) Delete(ctx context.Context, id uint64) error {
	oldItem, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	userResult, err := a.UserRepo.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		RoleIDs:         []uint64{id},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return errors.New400Response("不允许删除已经存在用户的角色")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		err := a.RoleMenuRepo.DeleteByRoleID(ctx, id)
		if err != nil {
			return err
		}

		return a.RoleRepo.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	a.Enforcer.DeleteRole(strconv.FormatUint(id, 10))

	return nil
}

// UpdateStatus 更新状态
func (a *RoleSrv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	oldItem, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Status == status {
		return nil
	}

	err = a.RoleRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	if status == 1 {
		roleMenus, err := a.RoleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{
			RoleID: id,
		})
		if err != nil {
			return err
		}

		resources, err := a.MenuActionResourceRepo.Query(ctx, schema.MenuActionResourceQueryParam{
			MenuIDs: roleMenus.Data.ToMenuIDs(),
		})
		if err != nil {
			return err
		}

		for _, ritem := range resources.Data.ToMap() {
			a.Enforcer.AddPermissionForUser(strconv.FormatUint(id, 10), ritem.Path, ritem.Method)
		}
	} else {
		a.Enforcer.DeleteRole(strconv.FormatUint(id, 10))
	}

	return nil
}
