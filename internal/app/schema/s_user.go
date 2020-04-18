package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/pkg/util"
)

// User 用户对象
type User struct {
	RecordID  string    `json:"record_id"`                             // 记录ID
	UserName  string    `json:"user_name" binding:"required"`          // 用户名
	RealName  string    `json:"real_name" binding:"required"`          // 真实姓名
	Password  string    `json:"password"`                              // 密码
	Phone     string    `json:"phone"`                                 // 手机号
	Email     string    `json:"email"`                                 // 邮箱
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 用户状态(1:启用 2:停用)
	Creator   string    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`    // 角色授权
}

func (a *User) String() string {
	return util.JSONMarshalToString(a)
}

// CleanSecure 清理安全数据
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// UserQueryParam 查询条件
type UserQueryParam struct {
	PaginationParam
	UserName   string   `form:"userName"`   // 用户名
	QueryValue string   `form:"queryValue"` // 模糊查询
	Status     int      `form:"status"`     // 用户状态(1:启用 2:停用)
	RoleIDs    []string `form:"-"`          // 角色ID列表
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserQueryResult 查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// ToShowResult 转换为显示结果
func (a UserQueryResult) ToShowResult(mUserRoles map[string]UserRoles, mRoles map[string]*Role) *UserShowQueryResult {
	return &UserShowQueryResult{
		PageResult: a.PageResult,
		Data:       a.Data.ToUserShows(mUserRoles, mRoles),
	}
}

// Users 用户对象列表
type Users []*User

// ToRecordIDs 转换为记录ID列表
func (a Users) ToRecordIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.RecordID
	}
	return idList
}

// ToUserShows 转换为用户显示列表
func (a Users) ToUserShows(mUserRoles map[string]UserRoles, mRoles map[string]*Role) UserShows {
	list := make(UserShows, len(a))
	for i, item := range a {
		showItem := new(UserShow)
		util.StructMapToStruct(item, showItem)
		for _, roleID := range mUserRoles[item.RecordID].ToRoleIDs() {
			if v, ok := mRoles[roleID]; ok {
				showItem.Roles = append(showItem.Roles, v)
			}
		}
		list[i] = showItem
	}

	return list
}

// ----------------------------------------UserRole--------------------------------------

// UserRole 用户角色
type UserRole struct {
	RecordID string `json:"record_id"` // 记录ID
	UserID   string `json:"user_id"`   // 用户ID
	RoleID   string `json:"role_id"`   // 角色ID
}

// UserRoleQueryParam 查询条件
type UserRoleQueryParam struct {
	PaginationParam
	UserID  string   // 用户ID
	UserIDs []string // 用户ID列表
}

// UserRoleQueryOptions 查询可选参数项
type UserRoleQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserRoleQueryResult 查询结果
type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *PaginationResult
}

// UserRoles 角色菜单列表
type UserRoles []*UserRole

// ToMap 转换为map
func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}

// ToRoleIDs 转换为角色ID列表
func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// ToUserIDMap 转换为用户ID映射
func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

// ----------------------------------------UserShow--------------------------------------

// UserShow 用户显示项
type UserShow struct {
	RecordID  string    `json:"record_id"`  // 记录ID
	UserName  string    `json:"user_name"`  // 用户名
	RealName  string    `json:"real_name"`  // 真实姓名
	Phone     string    `json:"phone"`      // 手机号
	Email     string    `json:"email"`      // 邮箱
	Status    int       `json:"status"`     // 用户状态(1:启用 2:停用)
	CreatedAt time.Time `json:"created_at"` // 创建时间
	Roles     []*Role   `json:"roles"`      // 授权角色列表
}

// UserShows 用户显示项列表
type UserShows []*UserShow

// UserShowQueryResult 用户显示项查询结果
type UserShowQueryResult struct {
	Data       UserShows
	PageResult *PaginationResult
}
