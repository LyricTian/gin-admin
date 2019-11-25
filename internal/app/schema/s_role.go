package schema

import "time"

// Role 角色对象
type Role struct {
	RecordID  string    `json:"record_id"`                     // 记录ID
	Name      string    `json:"name" binding:"required"`       // 角色名称
	Sequence  int       `json:"sequence"`                      // 排序值
	Memo      string    `json:"memo"`                          // 备注
	Creator   string    `json:"creator"`                       // 创建者
	CreatedAt time.Time `json:"created_at"`                    // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                    // 更新时间
	Menus     RoleMenus `json:"menus" binding:"required,gt=0"` // 菜单列表
}

// RoleQueryParam 查询条件
type RoleQueryParam struct {
	RecordIDs []string // 记录ID列表
	Name      string   // 角色名称
	LikeName  string   // 角色名称(模糊查询)
	UserID    string   // 用户ID
}

// RoleQueryOptions 查询可选参数项
type RoleQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
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
func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.RecordID] = item
	}
	return m
}
