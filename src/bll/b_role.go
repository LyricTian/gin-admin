package bll

import (
	"context"
	"time"

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
	UserModel model.IUser      `inject:"IUser"`
	Enforcer  *casbin.Enforcer `inject:""`
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
		return nil, errors.ErrNotFound
	}

	return item, nil
}

// 过滤叶子节点
func (a *Role) filterLeafMenuIDs(ctx context.Context, menuIDs []string) ([]string, error) {
	// menus, err := a.MenuModel.QuerySelect(ctx, schema.MenuSelectQueryParam{
	// 	RecordIDs: menuIDs,
	// 	Status:    1,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// var leafMenuIDs []string
	// for _, m := range menus {
	// 	var exists bool
	// 	for _, m2 := range menus {
	// 		if strings.HasPrefix(m2.LevelCode, m.LevelCode) &&
	// 			m2.LevelCode != m.LevelCode {
	// 			exists = true
	// 			break
	// 		}
	// 	}
	// 	if !exists {
	// 		leafMenuIDs = append(leafMenuIDs, m.RecordID)
	// 	}
	// }

	// return leafMenuIDs, nil
	return nil, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item *schema.Role) error {
	exists, err := a.RoleModel.CheckName(ctx, item.Name)
	if err != nil {
		return err
	} else if exists {
		return errors.NewBadRequestError("角色名称已经存在")
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
	err = a.RoleModel.Create(ctx, item)
	if err != nil {
		return err
	}

	return a.LoadPolicy(ctx, item.RecordID)
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item *schema.Role) error {
	oldItem, err := a.RoleModel.Get(ctx, recordID, false)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		exists, err := a.RoleModel.CheckName(ctx, item.Name)
		if err != nil {
			return err
		} else if exists {
			return errors.NewBadRequestError("角色名称已经存在")
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

	err = a.RoleModel.UpdateWithMenuIDs(ctx, recordID, info, item.MenuIDs)
	if err != nil {
		return err
	}

	return a.LoadPolicy(ctx, item.RecordID)
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	exists, err := a.RoleModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return errors.ErrNotFound
	}

	exists, err = a.UserModel.CheckByRoleID(ctx, recordID)
	if err != nil {
		return err
	} else if exists {
		return errors.NewBadRequestError("该角色已被赋予用户，不能删除！")
	}

	err = a.RoleModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}

	a.Enforcer.DeletePermissionsForUser(recordID)
	return nil
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, recordID string, status int) error {
	exists, err := a.RoleModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return errors.ErrNotFound
	}

	info := map[string]interface{}{
		"status": status,
	}

	err = a.RoleModel.Update(ctx, recordID, info)
	if err != nil {
		return err
	}

	if status == 2 {
		a.Enforcer.DeletePermissionsForUser(recordID)
	} else {
		err = a.LoadPolicy(ctx, recordID)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadAllPolicy 加载所有的角色策略
func (a *Role) LoadAllPolicy() error {
	ctx := context.Background()

	roles, err := a.RoleModel.QuerySelect(ctx, schema.RoleSelectQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	}

	for _, role := range roles {
		err = a.LoadPolicy(ctx, role.RecordID)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPolicy 加载角色权限策略
func (a *Role) LoadPolicy(ctx context.Context, roleID string) error {
	// menus, err := a.MenuModel.QuerySelect(ctx, schema.MenuSelectQueryParam{
	// 	Status: 1,
	// 	Types:  []int{40},
	// 	RoleID: roleID,
	// })
	// if err != nil {
	// 	return err
	// }

	// a.Enforcer.DeletePermissionsForUser(roleID)
	// for _, menu := range menus {
	// 	if menu.Path == "" || menu.Method == "" {
	// 		continue
	// 	}
	// 	a.Enforcer.AddPermissionForUser(roleID, menu.Path, menu.Method)
	// }

	return nil
}
