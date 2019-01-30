package schema

// Role 角色对象
type Role struct {
	RecordID string   `json:"record_id" swaggo:"false,记录ID"`
	Name     string   `json:"name" binding:"required" swaggo:"true,角色名称"`
	Sequence int      `json:"sequence" swaggo:"false,排序值"`
	Memo     string   `json:"memo" swaggo:"false,备注"`
	MenuIDs  []string `json:"menu_ids" binding:"required,gt=0" swaggo:"true,授权的菜单ID列表"`
}

// RoleQueryParam 角色对象查询条件
type RoleQueryParam struct {
	RecordIDs []string // 记录ID列表
	Name      string   // 角色名称
}

// RoleQueryOptions 角色对象查询可选参数项
type RoleQueryOptions struct {
	PageParam      *PaginationParam // 分页参数
	IncludeMenuIDs bool             // 是否包含菜单ID列表字段
}

// RoleQueryResult 角色对象查询结果
type RoleQueryResult struct {
	Data       Roles
	PageResult *PaginationResult
}

// RoleMini 角色对象
type RoleMini struct {
	RecordID string `json:"record_id" swaggo:"true,记录ID"`
	Name     string `json:"name" swaggo:"true,角色名称"`
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

// ToMiniList 转换为轻量数据
func (a Roles) ToMiniList() []*RoleMini {
	items := make([]*RoleMini, len(a))
	for i, item := range a {
		items[i] = &RoleMini{
			RecordID: item.RecordID,
			Name:     item.Name,
		}
	}
	return items
}
