package schema

// Menu 菜单管理
type Menu struct {
	ID        int64  `json:"id" db:"id,primarykey,autoincrement" structs:"id"`          // 唯一标识(自增ID)
	RecordID  string `json:"record_id" db:"record_id,size:36" structs:"record_id"`      // 记录内码(uuid)
	Code      string `json:"code" db:"code,size:50" structs:"code" binding:"required"`  // 菜单编号
	Name      string `json:"name" db:"name,size:50" structs:"name" binding:"required"`  // 菜单名称
	Type      int    `json:"type" db:"type" structs:"type" binding:"required"`          // 菜单类型(10：系统 20：模块 30：功能 40：资源)
	Sequence  int    `json:"sequence" db:"sequence" structs:"sequence"`                 // 排序值
	Icon      string `json:"icon" db:"icon,size:200" structs:"icon"`                    // 菜单图标
	Path      string `json:"path" db:"path,size:200" structs:"path"`                    // 访问路径
	Method    string `json:"method" db:"method,size:50" structs:"method"`               // 资源请求方式
	LevelCode string `json:"level_code" db:"level_code,size:20" structs:"level_code"`   // 分级码
	ParentID  string `json:"parent_id" db:"parent_id,size:36" structs:"parent_id"`      // 父级内码
	IsHide    int    `json:"is_hide" db:"is_hide" structs:"is_hide" binding:"required"` // 是否隐藏(1:是 2:否)
	Status    int    `json:"status" db:"status" structs:"status" binding:"required"`    // 状态(1:启用 2:停用)
	Creator   string `json:"creator" db:"creator,size:36" structs:"creator"`            // 创建人
	Created   int64  `json:"created" db:"created" structs:"created"`                    // 创建时间戳
	Updated   int64  `json:"updated" db:"updated" structs:"updated"`                    // 更新时间戳
	Deleted   int64  `json:"deleted" db:"deleted" structs:"deleted"`                    // 删除时间戳
}

// MenuQueryParam 菜单查询条件
type MenuQueryParam struct {
	Name     string // 菜单名称
	Type     int    // 菜单类型(10：系统 20：模块 30：功能 40：资源)
	ParentID string // 父级内码
	Status   int    // 状态(1:启用 2:停用)
}

// MenuQueryResult 菜单查询结果
type MenuQueryResult struct {
	ID       int64  `json:"id" db:"id"`               // 唯一标识(自增ID)
	RecordID string `json:"record_id" db:"record_id"` // 记录内码(uuid)
	Code     string `json:"code" db:"code"`           // 菜单编号
	Name     string `json:"name" db:"name"`           // 菜单名称
	Icon     string `json:"icon" db:"icon"`           // 菜单图标
	Path     string `json:"path" db:"path"`           // 访问路径
	Type     int    `json:"type" db:"type"`           // 菜单类型(10：系统 20：模块 30：功能 40：资源)
	Sequence int    `json:"sequence" db:"sequence"`   // 排序值
	IsHide   int    `json:"is_hide" db:"is_hide"`     // 是否隐藏(1:是 2:否)
	Status   int    `json:"status" db:"status"`       // 状态(1:启用 2:停用)
}

// MenuSelectQueryParam 菜单选择查询条件
type MenuSelectQueryParam struct {
	RecordIDs  []string // 记录ID列表
	Name       string   // 菜单名称
	Status     int      // 状态(1:启用 2:停用)
	UserID     string   // 用户ID
	RoleID     string   // 角色ID
	SystemCode string   // 系统编号
	IsHide     int      // 是否隐藏(1:是 2:否)
	Types      []int    // 菜单类型(10：系统 20：模块 30：功能 40：资源)
}

// MenuSelectQueryResult 菜单选择查询结果
type MenuSelectQueryResult struct {
	RecordID  string `json:"record_id" db:"record_id" structs:"record_id"`    // 记录内码(uuid)
	Code      string `json:"code" db:"code" structs:"code"`                   // 菜单编号
	Name      string `json:"name" db:"name" structs:"name"`                   // 菜单名称
	LevelCode string `json:"level_code" db:"level_code" structs:"level_code"` // 分级码
	ParentID  string `json:"parent_id" db:"parent_id" structs:"parent_id"`    // 父级内码
	Type      int    `json:"type" db:"type" structs:"type"`                   // 菜单类型(10：系统 20：模块 30：功能 40：资源)
	Icon      string `json:"icon" db:"icon" structs:"icon"`                   // 菜单图标
	Path      string `json:"path" db:"path" structs:"path"`                   // 访问路径
	Method    string `json:"method" db:"method" structs:"method"`             // 资源请求方式
}
