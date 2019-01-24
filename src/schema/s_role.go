package schema

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

// RoleMiniResult 角色对象最小数据结果
type RoleMiniResult struct {
	RecordID string `json:"record_id" swaggo:"true,记录ID"`
	Name     string `json:"name" swaggo:"true,角色名称"`
}

// Roles 角色对象列表
type Roles []*Role

// ToMap 转换为键值存储
func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.RecordID] = item
	}
	return m
}

// ToMiniResult 转换为角色最小结果
func (a Roles) ToMiniResult() []*RoleMiniResult {
	items := make([]*RoleMiniResult, len(a))
	for i, item := range items {
		items[i] = &RoleMiniResult{
			RecordID: item.RecordID,
			Name:     item.Name,
		}
	}
	return items
}
