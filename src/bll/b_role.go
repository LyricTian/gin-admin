package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/casbin/casbin"
)

// Role 角色管理
type Role struct {
	RoleModel model.IRole      `inject:"IRole"`
	MenuModel model.IMenu      `inject:"IMenu"`
	Enforcer  *casbin.Enforcer `inject:""`
	CommonBll *Common          `inject:""`
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
func (a *Role) QuerySelect(ctx context.Context) ([]*schema.RoleMini, error) {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{})
	if err != nil {
		return nil, err
	}
	return result.Data.ToMiniList(), nil
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string) (*schema.Role, error) {
	item, err := a.RoleModel.Get(ctx, recordID, schema.RoleQueryOptions{IncludeMenuIDs: true})
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Role) checkAndGetLeafMenuIDs(ctx context.Context, item schema.Role, oldItem *schema.Role) ([]string, error) {
	if oldItem == nil || oldItem.Name != item.Name {
		exists, err := a.RoleModel.CheckName(ctx, item.Name)
		if err != nil {
			return nil, err
		} else if exists {
			return nil, errors.NewBadRequestError("角色名称已经存在")
		}
	}

	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: item.MenuIDs,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.NewBadRequestError("请选择授权菜单")
	}

	return result.Data.ToLeafRecordIDs(), nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.Role, error) {
	leafMenuIDs, err := a.checkAndGetLeafMenuIDs(ctx, item, nil)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.MustUUID()
	item.MenuIDs = leafMenuIDs
	err = a.RoleModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicyWithRecordID(ctx, item.RecordID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item schema.Role) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	leafMenuIDs, err := a.checkAndGetLeafMenuIDs(ctx, item, nil)
	if err != nil {
		return err
	}
	item.MenuIDs = leafMenuIDs

	err = a.RoleModel.Update(ctx, recordID, item)
	if err != nil {
		return err
	}

	return a.LoadPolicyWithRecordID(ctx, recordID)
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordIDs ...string) error {
	err := a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		for _, recordID := range recordIDs {
			err := a.RoleModel.Delete(ctx, recordID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, recordID := range recordIDs {
		a.Enforcer.DeletePermissionsForUser(recordID)
	}
	return nil
}

// LoadAllPolicy 加载所有的角色策略
func (a *Role) LoadAllPolicy(ctx context.Context) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{},
		schema.RoleQueryOptions{IncludeMenuIDs: true})
	if err != nil {
		return err
	}

	for _, role := range result.Data {
		err = a.LoadPolicy(ctx, *role)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPolicyWithRecordID 加载角色权限策略
func (a *Role) LoadPolicyWithRecordID(ctx context.Context, recordID string) error {
	role, err := a.RoleModel.Get(ctx, recordID, schema.RoleQueryOptions{IncludeMenuIDs: true})
	if err != nil {
		return err
	} else if role == nil {
		return nil
	}

	return a.LoadPolicy(ctx, *role)
}

// LoadPolicy 加载角色权限策略
func (a *Role) LoadPolicy(ctx context.Context, item schema.Role) error {
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: item.MenuIDs,
		Types:     []int{3},
	})
	if err != nil {
		return err
	}

	roleID := item.RecordID
	a.Enforcer.DeletePermissionsForUser(roleID)
	for _, menu := range result.Data {
		if menu.Path == "" || menu.Method == "" {
			continue
		}
		a.Enforcer.AddPermissionForUser(roleID, menu.Path, menu.Method)
	}

	return nil
}
