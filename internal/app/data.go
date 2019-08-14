package app

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.uber.org/dig"
)

// InitData 初始化应用数据
func InitData(ctx context.Context, container *dig.Container) error {
	err := loadCasbinPolicyData(ctx, container)
	if err != nil {
		return err
	}

	if config.GetGlobalConfig().AllowInitMenu {
		return initMenuData(ctx, container)
	}

	return nil
}

// 加载casbin权限策略数据
func loadCasbinPolicyData(ctx context.Context, container *dig.Container) error {
	return container.Invoke(func(role bll.IRole, user bll.IUser) error {
		// 加载角色策略
		roleResult, err := role.Query(ctx, schema.RoleQueryParam{}, schema.RoleQueryOptions{
			IncludeMenus: true,
		})
		if err != nil {
			return err
		}

		for _, roleItem := range roleResult.Data {
			err := role.LoadPolicy(ctx, *roleItem)
			if err != nil {
				return err
			}
		}

		// 加载用户策略
		userResult, err := user.Query(ctx, schema.UserQueryParam{
			Status: 1,
		}, schema.UserQueryOptions{IncludeRoles: true})
		if err != nil {
			return err
		}

		for _, userItem := range userResult.Data {
			err := user.LoadPolicy(ctx, *userItem)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// initMenuData 初始化菜单数据
func initMenuData(ctx context.Context, container *dig.Container) error {
	return container.Invoke(func(trans bll.ITrans, menu bll.IMenu) error {
		// 检查是否存在菜单数据，如果不存在则初始化
		menuResult, err := menu.Query(ctx, schema.MenuQueryParam{}, schema.MenuQueryOptions{
			PageParam: &schema.PaginationParam{PageIndex: -1},
		})
		if err != nil {
			return err
		} else if menuResult.PageResult.Total > 0 {
			return nil
		}

		var data schema.MenuTrees
		err = util.JSONUnmarshal([]byte(menuData), &data)
		if err != nil {
			return err
		}

		return createMenus(ctx, trans, menu, "", data)
	})
}

func createMenus(ctx context.Context, trans bll.ITrans, menu bll.IMenu, parentID string, list schema.MenuTrees) error {
	return trans.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			sitem := schema.Menu{
				Name:      item.Name,
				Sequence:  item.Sequence,
				Icon:      item.Icon,
				Router:    item.Router,
				Hidden:    item.Hidden,
				ParentID:  parentID,
				Actions:   item.Actions,
				Resources: item.Resources,
			}
			nsitem, err := menu.Create(ctx, sitem)
			if err != nil {
				return err
			}

			if item.Children != nil && len(*item.Children) > 0 {
				err := createMenus(ctx, trans, menu, nsitem.RecordID, *item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// 初始化菜单数据
const menuData = `
[
  {
    "name": "首页",
    "icon": "dashboard",
    "router": "/dashboard",
    "sequence": 1900000
  },
  {
    "name": "系统管理",
    "icon": "setting",
    "sequence": 1100000,
    "children": [
      {
        "name": "菜单管理",
        "icon": "solution",
        "router": "/system/menu",
        "sequence": 1190000,
        "actions": [
          { "code": "add", "name": "新增" },
          { "code": "edit", "name": "编辑" },
          { "code": "del", "name": "删除" },
          { "code": "query", "name": "查询" }
        ],
        "resources": [
          {
            "code": "query",
            "name": "查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus"
          },
          {
            "code": "get",
            "name": "精确查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus/:id"
          },
          {
            "code": "create",
            "name": "创建菜单数据",
            "method": "POST",
            "path": "/api/v1/menus"
          },
          {
            "code": "update",
            "name": "更新菜单数据",
            "method": "PUT",
            "path": "/api/v1/menus/:id"
          },
          {
            "code": "delete",
            "name": "删除菜单数据",
            "method": "DELETE",
            "path": "/api/v1/menus/:id"
          }
        ]
      },
      {
        "name": "角色管理",
        "icon": "audit",
        "router": "/system/role",
        "sequence": 1180000,
        "actions": [
          { "code": "add", "name": "新增" },
          { "code": "edit", "name": "编辑" },
          { "code": "del", "name": "删除" },
          { "code": "query", "name": "查询" }
        ],
        "resources": [
          {
            "code": "query",
            "name": "查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles"
          },
          {
            "code": "get",
            "name": "精确查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles/:id"
          },
          {
            "code": "create",
            "name": "创建角色数据",
            "method": "POST",
            "path": "/api/v1/roles"
          },
          {
            "code": "update",
            "name": "更新角色数据",
            "method": "PUT",
            "path": "/api/v1/roles/:id"
          },
          {
            "code": "delete",
            "name": "删除角色数据",
            "method": "DELETE",
            "path": "/api/v1/roles/:id"
          },
          {
            "code": "queryMenu",
            "name": "查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus"
          }
        ]
      },
      {
        "name": "用户管理",
        "icon": "user",
        "router": "/system/user",
        "sequence": 1170000,
        "actions": [
          { "code": "add", "name": "新增" },
          { "code": "edit", "name": "编辑" },
          { "code": "del", "name": "删除" },
          { "code": "query", "name": "查询" }
        ],
        "resources": [
          {
            "code": "query",
            "name": "查询用户数据",
            "method": "GET",
            "path": "/api/v1/users"
          },
          {
            "code": "get",
            "name": "精确查询用户数据",
            "method": "GET",
            "path": "/api/v1/users/:id"
          },
          {
            "code": "create",
            "name": "创建用户数据",
            "method": "POST",
            "path": "/api/v1/users"
          },
          {
            "code": "update",
            "name": "更新用户数据",
            "method": "PUT",
            "path": "/api/v1/users/:id"
          },
          {
            "code": "delete",
            "name": "删除用户数据",
            "method": "DELETE",
            "path": "/api/v1/users/:id"
          },
          {
            "code": "disable",
            "name": "禁用用户数据",
            "method": "PATCH",
            "path": "/api/v1/users/:id/disable"
          },
          {
            "code": "enable",
            "name": "启用用户数据",
            "method": "PATCH",
            "path": "/api/v1/users/:id/enable"
          },
          {
            "code": "queryRole",
            "name": "查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles"
          }
        ]
      }
    ]
  }
]
`
