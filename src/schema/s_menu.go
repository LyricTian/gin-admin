package schema

import (
	"sort"
	"strings"

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

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	RecordIDs []string // 记录ID列表
	Code      string   // 菜单编号(模糊查询)
	Name      string   // 菜单名称(模糊查询)
	LevelCode string   // 分级码(前缀模糊查询)
	Types     []int    // 菜单类型(1：模块 2：功能 3：资源)
	IsHide    int      // 是否隐藏(1:是 2:否)
	ParentID  *string  // 父级内码
}

// MenuTreeResult 菜单树
type MenuTreeResult struct {
	RecordID string             `json:"record_id"`          // 记录内码
	Name     string             `json:"name"`               // 菜单名称
	ParentID string             `json:"parent_id"`          // 父级内码
	Children *[]*MenuTreeResult `json:"children,omitempty"` // 子级树
}

// Menus 菜单列表
type Menus []*Menu

func (a Menus) String() string {
	return util.JSONMarshalToString(a)
}

// ToLevelCodes 获取分级码列表（按照分级码正序排序）
func (a Menus) ToLevelCodes() []string {
	levelCodes := make([]string, len(a))
	for i, item := range a {
		levelCodes[i] = item.LevelCode
	}
	sort.Strings(levelCodes)
	return levelCodes
}

// ToTreeResult 转换为菜单树
func (a Menus) ToTreeResult() []*MenuTreeResult {
	var result []*MenuTreeResult
	_ = util.FillStructs(a, &result)

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
