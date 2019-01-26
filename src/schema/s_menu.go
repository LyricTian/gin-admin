package schema

import (
	"sort"
	"strings"
)

// Menu 菜单对象
type Menu struct {
	RecordID  string `json:"record_id" swaggo:"false,记录ID"`
	Code      string `json:"code" binding:"required" swaggo:"true,菜单编号"`
	Name      string `json:"name" binding:"required" swaggo:"true,菜单名称"`
	Type      int    `json:"type" binding:"required,max=3,min=1" swaggo:"true,菜单类型(1：模块 2：功能 3：资源)"`
	Sequence  int    `json:"sequence" swaggo:"false,排序值"`
	Icon      string `json:"icon" swaggo:"false,菜单图标"`
	Path      string `json:"path" swaggo:"false,访问路径"`
	Method    string `json:"method" swaggo:"false,资源请求方式"`
	LevelCode string `json:"level_code" swaggo:"false,分级码"`
	ParentID  string `json:"parent_id" swaggo:"false,父级内码"`
	IsHide    int    `json:"is_hide" binding:"required,max=2,min=1" swaggo:"true,是否隐藏(1:是 2:否)"`
}

// MenuQueryParam 菜单对象查询条件
type MenuQueryParam struct {
	RecordIDs  []string // 记录ID列表
	Code       string   // 菜单编号(模糊查询)
	Name       string   // 菜单名称(模糊查询)
	LevelCode  string   // 分级码(前缀模糊查询)
	Types      []int    // 菜单类型(1：模块 2：功能 3：资源)
	IsHide     int      // 是否隐藏(1:是 2:否)
	ParentID   *string  // 父级内码
	UserID     string   // 用户ID（查询用户所拥有的菜单权限）
	LevelCodes []string // 分级码列表
}

// MenuQueryOptions 菜单对象查询可选参数项
type MenuQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// MenuQueryResult 菜单对象查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// MenuTree 菜单树
type MenuTree struct {
	RecordID  string       `json:"record_id" swaggo:"false,记录ID"`
	Code      string       `json:"code" binding:"required" swaggo:"true,菜单编号"`
	Name      string       `json:"name" binding:"required" swaggo:"true,菜单名称"`
	Type      int          `json:"type" binding:"required,max=3,min=1" swaggo:"true,菜单类型(1：模块 2：功能 3：资源)"`
	Icon      string       `json:"icon" swaggo:"false,菜单图标"`
	Path      string       `json:"path" swaggo:"false,访问路径"`
	LevelCode string       `json:"level_code" swaggo:"false,分级码"`
	ParentID  string       `json:"parent_id" swaggo:"false,父级内码"`
	Children  *[]*MenuTree `json:"children,omitempty" swaggo:"false,子级树"`
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

// ToTree 转换为树形结构
func (a MenuTrees) ToTree() []*MenuTree {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.RecordID] = item
	}
	var data []*MenuTree
	for _, item := range a {
		if item.ParentID == "" {
			data = append(data, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				var children []*MenuTree
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return data
}

// Menus 菜单列表
type Menus []*Menu

// ToLevelCodes 获取分级码列表（按照分级码正序排序）
func (a Menus) ToLevelCodes() []string {
	levelCodes := make([]string, len(a))
	for i, item := range a {
		levelCodes[i] = item.LevelCode
	}
	sort.Strings(levelCodes)
	return levelCodes
}

// ToTreeList 转换为菜单树列表
func (a Menus) ToTreeList() MenuTrees {
	items := make([]*MenuTree, len(a))
	for i, item := range a {
		items[i] = &MenuTree{
			RecordID:  item.RecordID,
			Code:      item.Code,
			Name:      item.Name,
			Type:      item.Type,
			Icon:      item.Icon,
			Path:      item.Path,
			LevelCode: item.LevelCode,
			ParentID:  item.ParentID,
		}
	}
	return items
}

// ToLeafRecordIDs 转换为叶子节点记录ID列表
func (a Menus) ToLeafRecordIDs() []string {
	var recordIDs []string
	for _, item := range a {
		var exists bool
		for _, item2 := range a {
			if strings.HasPrefix(item2.LevelCode, item.LevelCode) {
				exists = false
				break
			}
		}
		if !exists {
			recordIDs = append(recordIDs, item.RecordID)
		}
	}
	return recordIDs
}
