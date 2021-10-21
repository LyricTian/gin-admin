package schema

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/pkg/util/hash"
	"github.com/LyricTian/gin-admin/v8/pkg/util/json"
	"github.com/LyricTian/gin-admin/v8/pkg/util/structure"
)

// GetRootUser 获取root用户
func GetRootUser() *User {
	user := config.C.Root
	return &User{
		ID:       user.UserID,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: hash.MD5String(user.Password),
	}
}

// CheckIsRootUser 检查是否是root用户
func CheckIsRootUser(ctx context.Context, userID uint64) bool {
	return GetRootUser().ID == userID
}

// User 用户对象
type User struct {
	ID        uint64    `json:"id,string"`                             // 唯一标识
	UserName  string    `json:"user_name" binding:"required"`          // 用户名
	RealName  string    `json:"real_name" binding:"required"`          // 真实姓名
	Password  string    `json:"password"`                              // 密码
	Phone     string    `json:"phone"`                                 // 手机号
	Email     string    `json:"email"`                                 // 邮箱
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 用户状态(1:启用 2:停用)
	Creator   uint64    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`    // 角色授权
}

func (a *User) String() string {
	return json.MarshalToString(a)
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
	RoleIDs    []uint64 `form:"-"`          // 角色ID列表
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	OrderFields  []*OrderField
	SelectFields []string
}

// UserQueryResult 查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// ToShowResult 转换为显示结果
func (a UserQueryResult) ToShowResult(mUserRoles map[uint64]UserRoles, mRoles map[uint64]*Role) *UserShowQueryResult {
	return &UserShowQueryResult{
		PageResult: a.PageResult,
		Data:       a.Data.ToUserShows(mUserRoles, mRoles),
	}
}

// Users 用户对象列表
type Users []*User

// ToIDs 转换为唯一标识列表
func (a Users) ToIDs() []uint64 {
	idList := make([]uint64, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// ToUserShows 转换为用户显示列表
func (a Users) ToUserShows(mUserRoles map[uint64]UserRoles, mRoles map[uint64]*Role) UserShows {
	list := make(UserShows, len(a))
	for i, item := range a {
		showItem := new(UserShow)
		structure.Copy(item, showItem)
		for _, roleID := range mUserRoles[item.ID].ToRoleIDs() {
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
	ID     uint64 `json:"id,string"`      // 唯一标识
	UserID uint64 `json:"user_id,string"` // 用户ID
	RoleID uint64 `json:"role_id,string"` // 角色ID
}

// UserRoleQueryParam 查询条件
type UserRoleQueryParam struct {
	PaginationParam
	UserID  uint64   // 用户ID
	UserIDs []uint64 // 用户ID列表
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
func (a UserRoles) ToMap() map[uint64]*UserRole {
	m := make(map[uint64]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}

// ToRoleIDs 转换为角色ID列表
func (a UserRoles) ToRoleIDs() []uint64 {
	list := make([]uint64, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// ToUserIDMap 转换为用户ID映射
func (a UserRoles) ToUserIDMap() map[uint64]UserRoles {
	m := make(map[uint64]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

// ----------------------------------------UserShow--------------------------------------

// UserShow 用户显示项
type UserShow struct {
	ID        uint64    `json:"id,string"`  // 唯一标识
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
