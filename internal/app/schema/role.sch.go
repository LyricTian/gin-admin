package schema

import (
	"fmt"
	"time"
)

// Role 角色对象
type Role struct {
	ID        uint64    `json:"id,string"`                             // 唯一标识
	Name      string    `json:"name" binding:"required"`               // 角色名称
	Sequence  int       `json:"sequence"`                              // 排序值
	Memo      string    `json:"memo"`                                  // 备注
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 状态(1:启用 2:禁用)
	Creator   uint64    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                            // 更新时间
	RoleMenus RoleMenus `json:"role_menus" binding:"required,gt=0"`    // 角色菜单列表
}

// RoleQueryParam 查询条件
type RoleQueryParam struct {
	PaginationParam
	IDs        []uint64 `form:"-"`          // 唯一标识列表
	Name       string   `form:"-"`          // 角色名称
	QueryValue string   `form:"queryValue"` // 模糊查询
	Status     int      `form:"status"`     // 状态(1:启用 2:禁用)
}

// RoleQueryOptions 查询可选参数项
type RoleQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// RoleQueryResult 查询结果
type RoleQueryResult struct {
	Data       Roles
	PageResult *PaginationResult
}

// Roles 角色对象列表
type Roles []*Role

// ToNames 获取角色名称列表
func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}
	return names
}

// ToMap 转换为键值存储
func (a Roles) ToMap() map[uint64]*Role {
	m := make(map[uint64]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// ----------------------------------------RoleMenu--------------------------------------

// RoleMenu 角色菜单对象
type RoleMenu struct {
	ID       uint64 `json:"id,string"`                           // 唯一标识
	RoleID   uint64 `json:"role_id,string" binding:"required"`   // 角色ID
	MenuID   uint64 `json:"menu_id,string" binding:"required"`   // 菜单ID
	ActionID uint64 `json:"action_id,string" binding:"required"` // 动作ID
}

// RoleMenuQueryParam 查询条件
type RoleMenuQueryParam struct {
	PaginationParam
	RoleID  uint64   // 角色ID
	RoleIDs []uint64 // 角色ID列表
}

// RoleMenuQueryOptions 查询可选参数项
type RoleMenuQueryOptions struct {
	OrderFields  []*OrderField
	SelectFields []string
}

// RoleMenuQueryResult 查询结果
type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *PaginationResult
}

// RoleMenus 角色菜单列表
type RoleMenus []*RoleMenu

// ToMap 转换为map
func (a RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range a {
		m[fmt.Sprintf("%d-%d", item.MenuID, item.ActionID)] = item
	}
	return m
}

// ToRoleIDMap 转换为角色ID映射
func (a RoleMenus) ToRoleIDMap() map[uint64]RoleMenus {
	m := make(map[uint64]RoleMenus)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}
	return m
}

// ToMenuIDs 转换为菜单ID列表
func (a RoleMenus) ToMenuIDs() []uint64 {
	var idList []uint64
	m := make(map[uint64]struct{})

	for _, item := range a {
		if _, ok := m[item.MenuID]; ok {
			continue
		}
		idList = append(idList, item.MenuID)
		m[item.MenuID] = struct{}{}
	}

	return idList
}

// ToActionIDs 转换为动作ID列表
func (a RoleMenus) ToActionIDs() []uint64 {
	idList := make([]uint64, len(a))
	m := make(map[uint64]struct{})
	for i, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}
	return idList
}
