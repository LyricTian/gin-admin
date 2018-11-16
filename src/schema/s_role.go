package schema

// Role 角色管理
type Role struct {
	ID       int64    `json:"id" db:"id,primarykey,autoincrement" structs:"id"`         // 唯一标识(自增ID)
	RecordID string   `json:"record_id" db:"record_id,size:36" structs:"record_id"`     // 记录内码(uuid)
	Name     string   `json:"name" db:"name,size:50" structs:"name" binding:"required"` // 角色名称
	Memo     string   `json:"memo" db:"memo,size:1024" structs:"memo"`                  // 角色备注
	Status   int      `json:"status" db:"status" structs:"status" binding:"required"`   // 角色状态(1:启用 2:停用)
	Creator  string   `json:"creator" db:"creator,size:36" structs:"creator"`           // 创建者
	Created  int64    `json:"created" db:"created" structs:"created"`                   // 创建时间戳
	Updated  int64    `json:"updated" db:"updated" structs:"updated"`                   // 更新时间戳
	Deleted  int64    `json:"deleted" db:"deleted" structs:"deleted"`                   // 删除时间戳
	MenuIDs  []string `json:"menu_ids" db:"-" structs:"-" binding:"required,gt=0"`      // 菜单ID列表
}

// RoleMenu 角色菜单管理
type RoleMenu struct {
	ID      int64  `json:"id" db:"id,primarykey,autoincrement"` // 唯一标识(自增ID)
	RoleID  string `json:"role_id" db:"role_id,size:36"`        // 角色内码
	MenuID  string `json:"menu_id" db:"menu_id,size:36"`        // 菜单内码
	Deleted int64  `json:"deleted" db:"deleted"`                // 删除时间戳
}

// RoleQueryParam 角色查询条件
type RoleQueryParam struct {
	Name   string // 角色名称
	Status int    // 角色状态(1:启用 2:停用)
}

// RoleQueryResult 角色查询结果
type RoleQueryResult struct {
	ID       int64  `json:"id" db:"id"`               // 唯一标识(自增ID)
	RecordID string `json:"record_id" db:"record_id"` // 记录内码
	Name     string `json:"name" db:"name"`           // 角色名称
	Memo     string `json:"memo" db:"memo"`           // 角色备注
	Status   int    `json:"status" db:"status"`       // 角色状态(1:启用 2:停用)
}

// RoleSelectQueryParam 角色选择查询条件
type RoleSelectQueryParam struct {
	RecordIDs []string // 记录ID列表
	Name      string   // 角色名称
	Status    int      // 状态(1:启用 2:停用)
}

// RoleSelectQueryResult 角色选择查询结果
type RoleSelectQueryResult struct {
	RecordID string `json:"record_id" db:"record_id"` // 记录内码
	Name     string `json:"name" db:"name"`           // 角色名称
}
