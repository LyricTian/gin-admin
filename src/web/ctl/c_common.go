package ctl

import (
	"context"
)

// Common API模块
type Common struct {
	DemoAPI  *Demo `inject:""`
	LoginAPI *Login
	UserAPI  *User
	RoleAPI  *Role
	MenuAPI  *Menu `inject:""`
}

// LoadCasbinPolicyData 加载casbin策略数据，包括角色权限数据、用户角色数据
func (c *Common) LoadCasbinPolicyData(ctx context.Context) error {
	// err := c.RoleAPI.RoleBll.LoadAllPolicy(ctx)
	// if err != nil {
	// 	return err
	// }

	// err = c.UserAPI.UserBll.LoadAllPolicy(ctx)
	// if err != nil {
	// 	return err
	// }
	return nil
}
