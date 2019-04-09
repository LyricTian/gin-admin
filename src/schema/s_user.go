package schema

import (
	"time"
)

// User 用户对象
type User struct {
	RecordID  string     `json:"record_id" swaggo:"false,记录ID"`
	UserName  string     `json:"user_name" binding:"required" swaggo:"true,用户名"`
	RealName  string     `json:"real_name" binding:"required" swaggo:"true,真实姓名"`
	Password  string     `json:"password" swaggo:"false,密码"`
	Phone     string     `json:"phone" swaggo:"false,手机号"`
	Email     string     `json:"email" swaggo:"false,邮箱"`
	Status    int        `json:"status" binding:"required,max=2,min=1" swaggo:"true,用户状态(1:启用 2:停用)"`
	Creator   string     `json:"creator" swaggo:"false,创建者"`
	CreatedAt *time.Time `json:"created_at" swaggo:"false,创建时间"`
	UpdatedAt *time.Time `json:"updated_at" swaggo:"false,更新时间"`
	Roles     UserRoles  `json:"roles" binding:"required,gt=0" swaggo:"true,角色授权"`
}

// UserRole 用户角色
type UserRole struct {
	RoleID string `json:"role_id" swaggo:"true,角色ID"`
}

// UserQueryParam 查询条件
type UserQueryParam struct {
	UserName     string   // 用户名
	LikeUserName string   // 用户名(模糊查询)
	LikeRealName string   // 真实姓名(模糊查询)
	Status       int      // 用户状态(1:启用 2:停用)
	RoleIDs      []string // 角色ID列表
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	PageParam    *PaginationParam // 分页参数
	IncludeRoles bool             // 包含角色权限
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
		roleIDs = append(roleIDs, item.Roles.ToRoleIDs()...)
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
			Email:     item.Email,
			Phone:     item.Phone,
			Status:    item.Status,
			CreatedAt: *item.CreatedAt,
		}

		var roles Roles
		for _, roleID := range item.Roles.ToRoleIDs() {
			if v, ok := mroles[roleID]; ok {
				roles = append(roles, v)
			}
		}
		show.Roles = roles
		list[i] = show
	}

	return list
}

// UserRoles 用户角色列表
type UserRoles []*UserRole

// ToRoleIDs 转换为角色ID列表
func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// UserPageShow 用户对象分页显示项
type UserPageShow struct {
	RecordID  string    `json:"record_id" swaggo:"false,记录ID"`
	UserName  string    `json:"user_name" swaggo:"true,用户名"`
	RealName  string    `json:"real_name" swaggo:"true,真实姓名"`
	Phone     string    `json:"phone" swaggo:"false,手机号"`
	Email     string    `json:"email" swaggo:"false,邮箱"`
	Status    int       `json:"status" swaggo:"true,用户状态(1:启用 2:停用)"`
	CreatedAt time.Time `json:"created_at" swaggo:"false,创建时间"`
	Roles     []*Role   `json:"roles" swaggo:"true,授权角色列表"`
}
