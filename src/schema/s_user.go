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
	Creator   string    `json:"creator" swaggo:"false,创建者"`
	CreatedAt time.Time `json:"created_at" swaggo:"false,创建时间"`
}

// UserQueryParam 查询条件
type UserQueryParam struct {
	UserName string   // 用户名(模糊查询)
	RealName string   // 真实姓名(模糊查询)
	Status   int      // 用户状态(1:启用 2:停用)
	RoleIDs  []string // 角色ID列表
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	PageParam      *PaginationParam // 分页参数
	IncludeRoleIDs bool             // 是否包含角色ID列表字段
}

// UserQueryResult 查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
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

// ToPageShows 转换为分页显示列表
func (a Users) ToPageShows(mroles map[string]*Role) []*UserPageShow {
	list := make([]*UserPageShow, len(a))

	for i, item := range a {
		show := &UserPageShow{
			RecordID:  item.RecordID,
			RealName:  item.RealName,
			UserName:  item.UserName,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		}

		var roles Roles
		for _, roleID := range item.RoleIDs {
			if v, ok := mroles[roleID]; ok {
				roles = append(roles, v)
			}
		}

		show.Roles = roles.ToMinis()
		list[i] = show
	}

	return list
}

// UserPageShow 用户对象分页显示项
type UserPageShow struct {
	RecordID  string      `json:"record_id" swaggo:"false,记录ID"`
	UserName  string      `json:"user_name" swaggo:"true,用户名"`
	RealName  string      `json:"real_name" swaggo:"true,真实姓名"`
	Status    int         `json:"status" swaggo:"true,用户状态(1:启用 2:停用)"`
	CreatedAt time.Time   `json:"created_at" swaggo:"false,创建时间"`
	Roles     []*RoleMini `json:"roles" swaggo:"true,授权角色列表"`
}
