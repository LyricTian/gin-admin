package ctl

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// Common 控制器公共模块
type Common struct {
	MenuCtl  *Menu  `inject:""`
	RoleCtl  *Role  `inject:""`
	UserCtl  *User  `inject:""`
	LoginCtl *Login `inject:""`
	DemoCtl  *Demo  `inject:""`
}

// LoadCasbinPolicyData 加载casbin策略数据，包括角色权限数据、用户角色数据
func (a *Common) LoadCasbinPolicyData(ctx context.Context) error {
	err := a.RoleCtl.RoleBll.LoadAllPolicy(ctx)
	if err != nil {
		return err
	}

	err = a.UserCtl.UserBll.LoadAllPolicy(ctx)
	if err != nil {
		return err
	}
	return nil
}

// InitMenuData 初始化菜单数据
func (a *Common) InitMenuData(ctx context.Context) error {
	// 检查是否存在菜单数据，如果不存在则初始化
	exists, err := a.MenuCtl.MenuBll.CheckDataInit(ctx)
	if err != nil {
		return err
	} else if exists {
		return nil
	}

	var data schema.MenuTrees
	err = util.JSONUnmarshal([]byte(menuData), &data)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, "", data)
}

func (a *Common) createMenus(ctx context.Context, parentID string, list schema.MenuTrees) error {
	return a.MenuCtl.MenuBll.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {

		for _, item := range list {
			sitem := schema.Menu{
				Name:      item.Name,
				Sequence:  item.Sequence,
				Icon:      item.Icon,
				Router:    item.Router,
				Hidden:    item.Hidden,
				ParentID:  parentID,
				Resources: item.Resources,
			}
			nsitem, err := a.MenuCtl.MenuBll.Create(ctx, sitem)
			if err != nil {
				return err
			}

			if item.Children != nil && len(*item.Children) > 0 {
				err := a.createMenus(ctx, nsitem.RecordID, *item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
