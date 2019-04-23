package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/casbin/casbin"
)

// NewRole 创建菜单管理实例
func NewRole(m *model.Common, e *casbin.Enforcer) *Role {
	return &Role{
		RoleModel: m.Role,
		MenuModel: m.Menu,
		Enforcer:  e,
	}
}

// Role 角色管理
type Role struct {
	RoleModel model.IRole
	MenuModel model.IMenu
	Enforcer  *casbin.Enforcer
}

// QueryPage 查询分页数据
func (a *Role) QueryPage(ctx context.Context, params schema.RoleQueryParam, pp *schema.PaginationParam) ([]*schema.Role, *schema.PaginationResult, error) {
	result, err := a.RoleModel.Query(ctx, params, schema.RoleQueryOptions{
		PageParam: pp,
	})
	if err != nil {
		return nil, nil, err
	}
	return result.Data, result.PageResult, nil
}

// QuerySelect 查询选择数据
func (a *Role) QuerySelect(ctx context.Context) ([]*schema.Role, error) {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{})
	if err != nil {
		return nil, err
	}

	// 清空部分字段数据
	return result.Data.ForEach(func(item *schema.Role, _ int) {
		item.Memo = ""
		item.Sequence = 0
		item.Creator = ""
		item.CreatedAt = nil
		item.UpdatedAt = nil
	}), nil
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string) (*schema.Role, error) {
	item, err := a.RoleModel.Get(ctx, recordID, schema.RoleQueryOptions{IncludeMenus: true})
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
		return errors.NewBadRequestError("角色名称已经存在")
	}
	return nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.Role, error) {
	err := a.checkName(ctx, item.Name)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.MustUUID()
	item.Creator = GetUserID(ctx)
	err = a.RoleModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	nitem, err := a.Get(ctx, item.RecordID)
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, nitem)
	if err != nil {
		return nil, err
	}
	return nitem, nil
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

	nitem, err := a.Get(ctx, item.RecordID)
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, nitem)
	if err != nil {
		return nil, err
	}
	return nitem, nil
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	// TODO: 如果用户已经被赋予该角色，则不允许删除
	err := a.RoleModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}

	a.Enforcer.DeletePermissionsForUser(recordID)
	return nil
}

// LoadAllPolicy 加载所有的角色策略
func (a *Role) LoadAllPolicy(ctx context.Context) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{},
		schema.RoleQueryOptions{IncludeMenus: true})
	if err != nil {
		return err
	}

	for _, role := range result.Data {
		err = a.LoadPolicy(ctx, role)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPolicyWithRecordID 加载角色权限策略
func (a *Role) LoadPolicyWithRecordID(ctx context.Context, recordID string) error {
	role, err := a.RoleModel.Get(ctx, recordID, schema.RoleQueryOptions{IncludeMenus: true})
	if err != nil {
		return err
	} else if role == nil {
		return nil
	}

	return a.LoadPolicy(ctx, role)
}

// LoadPolicy 加载角色权限策略
func (a *Role) LoadPolicy(ctx context.Context, item *schema.Role) error {
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
