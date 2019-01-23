package schema

import (
	"github.com/LyricTian/gin-admin/src/util"
)

// Role 角色对象
type Role struct {
	RecordID string   `json:"record_id" swaggo:"false,记录ID"`
	Name     string   `json:"name" binding:"required" swaggo:"true,角色名称"`
	Memo     string   `json:"memo" swaggo:"false,备注"`
	Status   int      `json:"status" binding:"required,max=2,min=1" swaggo:"true,状态(1:启用 2:停用)"`
	MenuIDs  []string `json:"menu_ids" binding:"required,gt=0" swaggo:"true,授权的菜单ID列表"`
}

// RoleQueryParam 角色查询条件
type RoleQueryParam struct {
	RecordIDs      []string // 记录ID列表
	Name           string   // 角色名称
	Status         int      // 角色状态(1:启用 2:停用)
	IncludeMenuIDs bool     // 是否包含菜单ID列表
}

// RoleSelectResult 角色选择结果
type RoleSelectResult struct {
	RecordID string `json:"record_id" swaggo:"true,记录ID"`
	Name     string `json:"name" swaggo:"true,角色名称"`
}

// Roles 角色对象列表
type Roles []*Role

// ToSelectResult 转换为角色选择结果列表
func (a Roles) ToSelectResult() []*RoleSelectResult {
	var items []*RoleSelectResult
	_ = util.FillStructs(a, &items)
	return items
}
