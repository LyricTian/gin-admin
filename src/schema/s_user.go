package schema

// User 用户管理
type User struct {
	ID       int64    `json:"id" db:"id,primarykey,autoincrement" structs:"id"`                        // 唯一标识(自增ID)
	RecordID string   `json:"record_id" db:"record_id,size:36" structs:"record_id"`                    // 记录内码(uuid)
	UserName string   `json:"user_name" db:"user_name,size:50" structs:"user_name" binding:"required"` // 用户名
	RealName string   `json:"real_name" db:"real_name,size:50" structs:"real_name" binding:"required"` // 真实姓名
	Password string   `json:"password" db:"password,size:40" structs:"password"`                       // 登录密码(sha1(md5(明文))加密)
	Status   int      `json:"status" db:"status" structs:"status" binding:"required"`                  // 用户状态(1:启用 2:停用)
	Creator  string   `json:"creator" db:"creator,size:36" structs:"creator"`                          // 创建者
	Created  int64    `json:"created" db:"created" structs:"created"`                                  // 创建时间戳
	Updated  int64    `json:"updated" db:"updated" structs:"updated"`                                  // 更新时间戳
	Deleted  int64    `json:"deleted" db:"deleted" structs:"deleted"`                                  // 删除时间戳
	RoleIDs  []string `json:"role_ids" db:"-" structs:"-" binding:"required,gt=0"`                     // 角色ID列表
}

// UserRole 用户角色授权管理
type UserRole struct {
	ID      int64  `json:"id" db:"id,primarykey,autoincrement"` // 唯一标识(自增ID)
	UserID  string `json:"user_id" db:"user_id,size:36"`        // 用户内码
	RoleID  string `json:"role_id" db:"role_id,size:36"`        // 角色内码
	Deleted int64  `json:"deleted" db:"deleted"`                // 删除时间戳
}

// UserQueryParam 用户查询条件
type UserQueryParam struct {
	UserName string // 用户名
	RealName string // 真实姓名
	Status   int    // 用户状态(1:启用 2:停用)
	RoleID   string // 角色ID
}

// UserQueryResult 用户查询结果
type UserQueryResult struct {
	ID        int64    `json:"id" db:"id"`               // 唯一标识(自增ID)
	RecordID  string   `json:"record_id" db:"record_id"` // 记录内码(uuid)
	UserName  string   `json:"user_name" db:"user_name"` // 用户名
	RealName  string   `json:"real_name" db:"real_name"` // 真实姓名
	Status    int      `json:"status" db:"status"`       // 用户状态(1:启用 2:停用)
	Created   int64    `json:"created" db:"created"`     // 创建时间戳
	RoleNames []string `json:"role_names" db:"-"`        // 角色名称
}

// UserRoleQueryParam 用户角色查询参数
type UserRoleQueryParam struct {
	UserID string
}
