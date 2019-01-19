package schema

import (
	"sort"

	"github.com/LyricTian/gin-admin/src/util"
)

// Menu 菜单管理
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

// MenuPageQueryParam 分页查询条件
type MenuPageQueryParam struct {
	Code     string // 菜单编号(模糊查询)
	Name     string // 菜单名称(模糊查询)
	Type     int    // 菜单类型(1：模块 2：功能 3：资源)
	ParentID string // 父级内码
}

// MenuListQueryParam 列表查询条件
type MenuListQueryParam struct {
	LevelCode  string   // 分级码(模糊查询)
	LevelCodes []string // 分级码列表
	Types      []int    // 菜单类型(1：模块 2：功能 3：资源)
	IsHide     int      // 是否隐藏(1:是 2:否)
	ParentID   *string  // 父级内码
}

// MenuTreeResult 菜单树
type MenuTreeResult struct {
	RecordID string             `json:"record_id"`          // 记录内码
	Name     string             `json:"name"`               // 菜单名称
	ParentID string             `json:"parent_id"`          // 父级内码
	Children *[]*MenuTreeResult `json:"children,omitempty"` // 子级树
}

// MenuList 菜单列表
type MenuList []*Menu

func (m MenuList) String() string {
	return util.JSONMarshalToString(m)
}

// ToLevelCodes 获取分级码列表（按照分级码正序排序）
func (m MenuList) ToLevelCodes() []string {
	levelCodes := make([]string, len(m))
	for i, item := range m {
		levelCodes[i] = item.LevelCode
	}
	sort.Strings(levelCodes)
	return levelCodes
}

// ToTreeResult 转换为菜单树
func (m MenuList) ToTreeResult() []*MenuTreeResult {
	var result []*MenuTreeResult
	_ = util.FillStructs(m, &result)

	mi := make(map[string]*MenuTreeResult)
	for _, item := range result {
		mi[item.RecordID] = item
	}

	var data []*MenuTreeResult
	for _, item := range result {
		if item.ParentID == "" {
			data = append(data, item)
			continue
		}

		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				var children []*MenuTreeResult
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}

	return data
}
