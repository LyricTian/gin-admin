package schema

import (
	"time"
)

// User 用户对象
type User struct {
	RecordID  string    `json:"record_id" swaggo:"false,记录ID"`
	UserName  string    `json:"user_name" binding:"required" swaggo:"true,用户名"`
	RealName  string    `json:"real_name" binding:"required" swaggo:"true,真实姓名"`
	Password  string    `json:"password" swaggo:"false,密码"`
	Status    int       `json:"status" binding:"required,max=2,min=1" swaggo:"true,用户状态(1:启用 2:停用)"`
	RoleIDs   []string  `json:"role_ids" binding:"required,gt=0" swaggo:"true,授权角色ID列表"`
	CreatedAt time.Time `json:"created_at" swaggo:"false,创建时间"`
}

// UserQueryParam 用户查询条件
type UserQueryParam struct {
	UserName string // 用户名
	RealName string // 真实姓名
	Status   int    // 用户状态(1:启用 2:停用)
	RoleID   string // 角色ID
}

// UserQueryOptions 用户对象查询可选参数项
type UserQueryOptions struct {
	PageParam       *PaginationParam // 分页参数
	IncludePassword bool             // 是否包含密码字段
	IncludeRoleIDs  bool             // 是否包含角色ID列表字段
}

// UserQueryResult 用户查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// UserPageQueryResult 用户分页查询结果
type UserPageQueryResult struct {
	RecordID  string            `json:"record_id" swaggo:"false,记录ID"`
	UserName  string            `json:"user_name" swaggo:"true,用户名"`
	RealName  string            `json:"real_name" swaggo:"true,真实姓名"`
	Status    int               `json:"status" swaggo:"true,用户状态(1:启用 2:停用)"`
	CreatedAt time.Time         `json:"created_at" swaggo:"false,创建时间"`
	Roles     []*RoleMiniResult `json:"roles" swaggo:"true,授权角色列表"`
}

// Users 用户对象列表
type Users []*User

// ToRoleIDs 获取角色ID列表
func (a Users) ToRoleIDs() []string {
	var roleIDs []string
	for _, item := range a {
		roleIDs = append(roleIDs, item.RoleIDs...)
	}
	return roleIDs
}

// ToPageQueryResult 转换为分页查询结果
func (a Users) ToPageQueryResult(roles map[string]*Role) []*UserPageQueryResult {
	items := make([]*UserPageQueryResult, len(a))

	for i, user := range a {
		result := &UserPageQueryResult{
			RecordID:  user.RecordID,
			RealName:  user.RealName,
			UserName:  user.UserName,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		}

		var roleItems []*Role
		for _, roleID := range user.RoleIDs {
			if v, ok := roles[roleID]; ok {
				roleItems = append(roleItems, v)
			}
		}
		result.Roles = Roles(roleItems).ToMiniResult()

		items[i] = result
	}

	return items
}
