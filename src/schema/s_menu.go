package schema

import (
	"sort"

	"github.com/LyricTian/gin-admin/src/util"
)

// Menu 菜单管理
type Menu struct {
	RecordID  string `json:"record_id"`                  // 记录内码
	Code      string `json:"code" binding:"required"`    // 菜单编号
	Name      string `json:"name" binding:"required"`    // 菜单名称
	Type      int    `json:"type" binding:"required"`    // 菜单类型(10：模块 20：功能 30：资源)
	Sequence  int    `json:"sequence"`                   // 排序值
	Icon      string `json:"icon"`                       // 菜单图标
	Path      string `json:"path"`                       // 访问路径
	Method    string `json:"method"`                     // 资源请求方式
	LevelCode string `json:"level_code"`                 // 分级码
	ParentID  string `json:"parent_id"`                  // 父级内码
	IsHide    int    `json:"is_hide" binding:"required"` // 是否隐藏(1:是 2:否)
}

// MenuPageQueryParam 分页查询条件
type MenuPageQueryParam struct {
	Code     string // 菜单编号(模糊查询)
	Name     string // 菜单名称(模糊查询)
	Type     int    // 菜单类型(10：模块 20：功能 30：资源)
	ParentID string // 父级内码
}

// MenuListQueryParam 列表查询条件
type MenuListQueryParam struct {
	LevelCode  string   // 分级码(模糊查询)
	LevelCodes []string // 分级码列表
	Types      []int    // 菜单类型(10：模块 20：功能 30：资源)
	IsHide     int      // 是否隐藏(1:是 2:否)
	ParentID   string   // 父级内码
}

// MenuTreeResult 菜单树
type MenuTreeResult struct {
	RecordID string            `json:"record_id"` // 记录内码
	Name     string            `json:"name"`      // 菜单名称
	ParentID string            `json:"parent_id"` // 父级内码
	Children []*MenuTreeResult `json:"children"`  // 子级树
}

// MenuList 菜单列表
type MenuList []*Menu

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

		if pitem, ok := mi[item.ParentID]; !ok {
			pitem.Children = append(pitem.Children, item)
			continue
		}

	}

	return data
}
