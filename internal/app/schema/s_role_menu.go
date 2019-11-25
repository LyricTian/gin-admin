package schema

// RoleMenu 角色菜单对象
type RoleMenu struct {
	RecordID string `json:"record_id"`                    // 记录ID
	RoleID   string `json:"role_id" binding:"required"`   // 角色ID
	MenuID   string `json:"menu_id" binding:"required"`   // 菜单ID
	ActionID string `json:"action_id" binding:"required"` // 动作ID
}

// RoleMenuQueryParam 查询条件
type RoleMenuQueryParam struct {
	RoleID string // 角色ID
}

// RoleMenuQueryOptions 查询可选参数项
type RoleMenuQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
}

// RoleMenuQueryResult 查询结果
type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *PaginationResult
}

// RoleMenus 角色菜单列表
type RoleMenus []*RoleMenu
