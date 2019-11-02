package app

import (
	"context"
	"os"

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

	if c := config.Global(); c.AllowInitMenu && c.Menu != "" {
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
			if roleItem != nil {
				err := role.LoadPolicy(ctx, *roleItem)
				if err != nil {
					return err
				}
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
			if userItem != nil {
				err := user.LoadPolicy(ctx, *userItem)
				if err != nil {
					return err
				}
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

		data, err := readMenuData()
		if err != nil {
			return err
		}

		return createMenus(ctx, trans, menu, "", data)
	})
}

func readMenuData() (schema.MenuTrees, error) {
	file, err := os.Open(config.Global().Menu)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data schema.MenuTrees
	err = util.JSONNewDecoder(file).Decode(&data)
	return data, err
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
