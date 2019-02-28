package schema

import "time"

// Role 角色对象
type Role struct {
	RecordID  string    `json:"record_id" swaggo:"false,记录ID"`
	Name      string    `json:"name" binding:"required" swaggo:"true,角色名称"`
	Sequence  int       `json:"sequence" swaggo:"false,排序值"`
	Memo      string    `json:"memo" swaggo:"false,备注"`
	MenuIDs   []string  `json:"menu_ids" binding:"required,gt=0" swaggo:"true,授权的菜单ID列表"`
	Creator   string    `json:"creator" swaggo:"false,创建者"`
	CreatedAt time.Time `json:"created_at" swaggo:"false,创建时间"`
}

// RoleQueryParam 查询条件
type RoleQueryParam struct {
	RecordIDs []string // 记录ID列表
	Name      string   // 角色名称(模糊查询)
}

// RoleQueryOptions 查询可选参数项
type RoleQueryOptions struct {
	PageParam      *PaginationParam // 分页参数
	IncludeMenuIDs bool             // 是否包含菜单ID列表字段
}

// RoleQueryResult 查询结果
type RoleQueryResult struct {
	Data       Roles
	PageResult *PaginationResult
}

// Roles 角色对象列表
type Roles []*Role

// ToNames 转换为名称列表
func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}
	return names
}

// ToMap 转换为键值存储
func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.RecordID] = item
	}
	return m
}

// ToMinis 转换为mini列表
func (a Roles) ToMinis() []*RoleMini {
	items := make([]*RoleMini, len(a))
	for i, item := range a {
		items[i] = &RoleMini{
			RecordID: item.RecordID,
			Name:     item.Name,
		}
	}
	return items
}

// RoleMini 角色对象(少量数据)
type RoleMini struct {
	RecordID string `json:"record_id" swaggo:"true,记录ID"`
	Name     string `json:"name" swaggo:"true,角色名称"`
}
