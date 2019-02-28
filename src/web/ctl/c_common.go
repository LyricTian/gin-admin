package ctl

import (
	"context"
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	wcontext "github.com/LyricTian/gin-admin/src/web/context"
)

// Common 控制器公共模块
type Common struct {
	DemoCtl     *Demo     `inject:""`
	LoginCtl    *Login    `inject:""`
	UserCtl     *User     `inject:""`
	RoleCtl     *Role     `inject:""`
	MenuCtl     *Menu     `inject:""`
	ResourceCtl *Resource `inject:""`
}

// LoadCasbinPolicyData 加载casbin策略数据，包括角色权限数据、用户角色数据
func (c *Common) LoadCasbinPolicyData(ctx context.Context) error {
	err := c.RoleCtl.RoleBll.LoadAllPolicy(ctx)
	if err != nil {
		return err
	}

	err = c.UserCtl.UserBll.LoadAllPolicy(ctx)
	if err != nil {
		return err
	}
	return nil
}

// CheckAndCreateResource 检查并创建资源数据
func (c *Common) CheckAndCreateResource(ctx context.Context) error {
	data := wcontext.GetRouterData()
	for k, item := range data {
		idx := strings.IndexByte(k, '/')
		if idx < 0 {
			continue
		}
		method, path := k[:idx], k[idx:]
		_, err := c.ResourceCtl.ResourceBll.Create(ctx, schema.Resource{
			Code:   item.Code,
			Name:   item.Name,
			Path:   path,
			Method: method,
		})
		if err != nil {
			if err == bll.ErrResourcePathAndMethodExists {
				continue
			}
			return err
		}
	}

	return nil
}

// InitMenuData 初始化菜单数据
func (c *Common) InitMenuData(ctx context.Context) error {
	data := `
[{
	"code":"dashboard",
	"name":"首页",
	"icon":"dashboard",
	"type":1,
	"sequence":9,
	"path":"/dashboard"
},{
	"code":"example",
	"name":"演示用例",
	"icon":"bulb",
	"type":1,
	"sequence":2,
	"children":[
		{
			"code":"demo",
			"name":"基础示例",
			"icon":"experiment",
			"type":2,
			"path":"/example/demo",
			"sequence":1
		}
	]
},{
	"code":"system",
	"name":"系统管理",
	"icon":"setting",
	"type":1,
	"sequence":1,
	"children":[
		{
			"code":"resource",
			"name":"资源管理",
			"icon":"api",
			"type":2,
			"sequence":9,
			"path":"/system/resource"
		},
		{
			"code":"menu",
			"name":"菜单管理",
			"icon":"solution",
			"type":2,
			"sequence":8,
			"path":"/system/menu"
		},
		{
			"code":"role",
			"name":"角色管理",
			"icon":"audit",
			"type":2,
			"sequence":7,
			"path":"/system/role"
		},
		{
			"code":"user",
			"name":"用户管理",
			"icon":"user",
			"type":2,
			"sequence":6,
			"path":"/system/user"
		}
	]
}]
`

	// 检查是否存在数据，如果存在则不执行初始化
	exists, err := c.MenuCtl.MenuBll.CheckDataInit(ctx)
	if err != nil {
		return err
	} else if exists {
		return nil
	}

	var items []*schema.MenuTree
	err = util.JSONUnmarshal([]byte(data), &items)
	if err != nil {
		return err
	}

	err = c.MenuCtl.MenuBll.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		return c.createMenu(ctx, "", items)
	})

	return err
}

func (c *Common) createMenu(ctx context.Context, parentID string, items []*schema.MenuTree) error {
	for _, item := range items {
		newItem, err := c.MenuCtl.MenuBll.Create(ctx, schema.Menu{
			Code:     item.Code,
			Name:     item.Name,
			Type:     item.Type,
			Sequence: item.Sequence,
			Icon:     item.Icon,
			Path:     item.Path,
			ParentID: parentID,
		})
		if err != nil {
			return err
		}

		if item.Children != nil {
			err = c.createMenu(ctx, newItem.RecordID, *item.Children)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
