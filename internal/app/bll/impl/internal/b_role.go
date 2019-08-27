package internal

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/casbin/casbin"
)

// NewRole 创建角色管理实例
func NewRole(
	e *casbin.Enforcer,
	mRole model.IRole,
	mMenu model.IMenu,
	mUser model.IUser,
) *Role {
	return &Role{
		Enforcer:  e,
		RoleModel: mRole,
		MenuModel: mMenu,
		UserModel: mUser,
	}
}

// Role 角色管理
type Role struct {
	Enforcer  *casbin.Enforcer
	RoleModel model.IRole
	MenuModel model.IMenu
	UserModel model.IUser
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

	return item, nil
}

func (a *Role) checkName(ctx context.Context, name string) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		Name: name,
	}, schema.RoleQueryOptions{
		PageParam: &schema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.ErrResourceExists
	}
	return nil
}

func (a *Role) getUpdate(ctx context.Context, recordID string) (*schema.Role, error) {
	nitem, err := a.Get(ctx, recordID, schema.RoleQueryOptions{
		IncludeMenus: true,
	})
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, *nitem)
	if err != nil {
		return nil, err
	}
	return nitem, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.Role, error) {
	err := a.checkName(ctx, item.Name)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.MustUUID()
	err = a.RoleModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, item.RecordID)
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item schema.Role) (*schema.Role, error) {
	oldItem, err := a.RoleModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		err := a.checkName(ctx, item.Name)
		if err != nil {
			return nil, err
		}
	}

	err = a.RoleModel.Update(ctx, recordID, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, recordID)
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	// 如果用户已经被赋予该角色，则不允许删除
	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		RoleIDs: []string{recordID},
	}, schema.UserQueryOptions{
		PageParam: &schema.PaginationParam{PageIndex: -1},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return errors.ErrResourceNotAllowDelete
	}

	err = a.RoleModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}

	a.Enforcer.DeletePermissionsForUser(recordID)
	return nil
}

// LoadPolicy 加载角色权限策略
func (a *Role) LoadPolicy(ctx context.Context, item schema.Role) error {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: item.Menus.ToMenuIDs(),
	}, schema.MenuQueryOptions{
		IncludeResources: true,
	})
	if err != nil {
		return err
	}

	menuMap := result.Data.ToMap()
	roleID := item.RecordID
	a.Enforcer.DeletePermissionsForUser(roleID)

	for _, item := range item.Menus {
		mitem, ok := menuMap[item.MenuID]
		if !ok {
			continue
		}
		resMap := mitem.Resources.ToMap()
		for _, res := range item.Resources {
			ritem, ok := resMap[res]
			if !ok || ritem.Path == "" || ritem.Method == "" {
				continue
			}
			a.Enforcer.AddPermissionForUser(roleID, ritem.Path, ritem.Method)
		}
	}

	return nil
}
