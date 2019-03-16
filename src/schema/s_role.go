package schema

import "time"

// Role 角色对象
type Role struct {
	RecordID  string     `json:"record_id" swaggo:"false,记录ID"`
	Name      string     `json:"name" binding:"required" swaggo:"true,角色名称"`
	Sequence  int        `json:"sequence" swaggo:"false,排序值"`
	Memo      string     `json:"memo" swaggo:"false,备注"`
	Creator   string     `json:"creator" swaggo:"false,创建者"`
	CreatedAt *time.Time `json:"created_at" swaggo:"false,创建时间"`
	UpdatedAt *time.Time `json:"updated_at" swaggo:"false,更新时间"`
	Menus     RoleMenus  `json:"menus" binding:"required,gt=0" swaggo:"false,菜单权限"`
}

// RoleMenu 角色菜单对象
type RoleMenu struct {
	MenuID    string   `json:"menu_id" swaggo:"false,菜单ID"`
	Actions   []string `json:"actions" swaggo:"false,动作权限列表"`
	Resources []string `json:"resources" swaggo:"false,资源权限列表"`
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
	PageParam    *PaginationParam // 分页参数
	IncludeMenus bool             // 包含菜单权限
}

// RoleQueryResult 查询结果
type RoleQueryResult struct {
	Data       Roles
	PageResult *PaginationResult
}

// Roles 角色对象列表
type Roles []*Role

// ForEach 遍历数据项
func (a Roles) ForEach(fn func(*Role, int)) Roles {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToMenuIDs 获取所有的菜单ID（不去重）
func (a Roles) ToMenuIDs() []string {
	var idList []string
	for _, item := range a {
		idList = append(idList, item.Menus.ToMenuIDs()...)
	}
	return idList
}

func (a Roles) mergeStrings(olds, news []string) []string {
	for _, n := range news {
		exists := false
		for _, o := range olds {
			if o == n {
				exists = true
				break
			}
		}
		if !exists {
			olds = append(olds, n)
		}
	}
	return olds
}

// ToMenuIDActionsMap 转换为菜单ID的动作权限列表映射
func (a Roles) ToMenuIDActionsMap() map[string][]string {
	m := make(map[string][]string)
	for _, item := range a {
		for _, menu := range item.Menus {
			v, ok := m[menu.MenuID]
			if ok {
				m[menu.MenuID] = a.mergeStrings(v, menu.Actions)
				continue
			}
			m[menu.MenuID] = menu.Actions
		}
	}
	return m
}

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

// RoleMenus 角色菜单列表
type RoleMenus []*RoleMenu

// ToMenuIDs 转换为菜单ID列表
func (a RoleMenus) ToMenuIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.MenuID
	}
	return list
}
