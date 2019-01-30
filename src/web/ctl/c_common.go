package ctl

import (
	"context"
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/schema"
	wcontext "github.com/LyricTian/gin-admin/src/web/context"
)

// Common API模块
type Common struct {
	DemoAPI     *Demo     `inject:""`
	LoginAPI    *Login    `inject:""`
	UserAPI     *User     `inject:""`
	RoleAPI     *Role     `inject:""`
	MenuAPI     *Menu     `inject:""`
	ResourceAPI *Resource `inject:""`
}

// LoadCasbinPolicyData 加载casbin策略数据，包括角色权限数据、用户角色数据
func (c *Common) LoadCasbinPolicyData(ctx context.Context) error {
	err := c.RoleAPI.RoleBll.LoadAllPolicy(ctx)
	if err != nil {
		return err
	}

	err = c.UserAPI.UserBll.LoadAllPolicy(ctx)
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
		method, path := k[:idx], k[idx+1:]
		_, err := c.ResourceAPI.ResourceBll.Create(ctx, schema.Resource{
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
