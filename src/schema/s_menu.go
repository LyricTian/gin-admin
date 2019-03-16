package schema

import (
	"strings"
	"time"
)

// Menu 菜单对象
type Menu struct {
	RecordID   string        `json:"record_id" swaggo:"false,记录ID"`
	Name       string        `json:"name" binding:"required" swaggo:"true,菜单名称"`
	Sequence   int           `json:"sequence" swaggo:"false,排序值"`
	Icon       string        `json:"icon" swaggo:"false,菜单图标"`
	Router     string        `json:"router" swaggo:"false,访问路由"`
	Hidden     int           `json:"hidden" swaggo:"false,隐藏菜单(0:不隐藏 1:隐藏)"`
	ParentID   string        `json:"parent_id" swaggo:"false,父级ID"`
	ParentPath string        `json:"parent_path" swaggo:"false,父级路径"`
	Creator    string        `json:"creator" swaggo:"false,创建者"`
	CreatedAt  *time.Time    `json:"created_at" swaggo:"false,创建时间"`
	UpdatedAt  *time.Time    `json:"updated_at" swaggo:"false,更新时间"`
	Actions    MenuActions   `json:"actions" swaggo:"false,动作列表"`
	Resources  MenuResources `json:"resources" swaggo:"false,资源列表"`
}

// MenuAction 菜单动作对象
type MenuAction struct {
	Code string `json:"code" swaggo:"true,动作编号"`
	Name string `json:"name" swaggo:"true,动作名称"`
}

// MenuResource 菜单资源对象
type MenuResource struct {
	Code   string `json:"code" swaggo:"true,资源编号"`
	Name   string `json:"name" swaggo:"true,资源名称"`
	Method string `json:"method" swaggo:"true,请求方式"`
	Path   string `json:"path" swaggo:"true,请求路径"`
}

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	RecordIDs        []string // 记录ID列表
	LikeName         string   // 菜单名称(模糊查询)
	ParentID         *string  // 父级内码
	PrefixParentPath string   // 父级路径(前缀模糊查询)
	Hidden           *int     // 隐藏菜单
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	PageParam        *PaginationParam // 分页参数
	IncludeActions   bool             // 包含动作列表
	IncludeResources bool             // 包含资源列表
}

// MenuQueryResult 查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// Menus 菜单列表
type Menus []*Menu

// ToMap 转换为键值映射
func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.RecordID] = item
	}
	return m
}

// SplitAndGetAllRecordIDs 拆分父级路径并获取所有记录ID
func (a Menus) SplitAndGetAllRecordIDs() []string {
	var recordIDs []string
	for _, item := range a {
		recordIDs = append(recordIDs, item.RecordID)
		if item.ParentPath == "" {
			continue
		}

		pps := strings.Split(item.ParentPath, "/")
		for _, pp := range pps {
			var exists bool
			for _, recordID := range recordIDs {
				if pp == recordID {
					exists = true
					break
				}
			}
			if !exists {
				recordIDs = append(recordIDs, pp)
			}
		}
	}
	return recordIDs
}

// ToTrees 转换为菜单列表
func (a Menus) ToTrees() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			RecordID:   item.RecordID,
			Name:       item.Name,
			Sequence:   item.Sequence,
			Icon:       item.Icon,
			Router:     item.Router,
			Hidden:     item.Hidden,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Actions:    item.Actions,
			Resources:  item.Resources,
		}
	}
	return list
}

func (a Menus) fillLeafNodeID(tree *[]*MenuTree, leafNodeIDs *[]string) {
	for _, node := range *tree {
		if node.Children == nil || len(*node.Children) == 0 {
			*leafNodeIDs = append(*leafNodeIDs, node.RecordID)
			continue
		}
		a.fillLeafNodeID(node.Children, leafNodeIDs)
	}
}

// ToLeafRecordIDs 转换为叶子节点记录ID列表
func (a Menus) ToLeafRecordIDs() []string {
	var leafNodeIDs []string
	tree := a.ToTrees().ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

// MenuResources 菜单资源列表
type MenuResources []*MenuResource

// ForEach 遍历资源数据
func (a MenuResources) ForEach(fn func(*MenuResource, int)) MenuResources {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToMap 转换为键值映射
func (a MenuResources) ToMap() map[string]*MenuResource {
	m := make(map[string]*MenuResource)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

// MenuActions 菜单动作列表
type MenuActions []*MenuAction

// MenuTree 菜单树
type MenuTree struct {
	RecordID   string        `json:"record_id" swaggo:"false,记录ID"`
	Name       string        `json:"name" binding:"required" swaggo:"true,菜单名称"`
	Sequence   int           `json:"sequence" swaggo:"false,排序值"`
	Icon       string        `json:"icon" swaggo:"false,菜单图标"`
	Router     string        `json:"router" swaggo:"false,访问路由"`
	Hidden     int           `json:"hidden" swaggo:"false,隐藏菜单(0:不隐藏 1:隐藏)"`
	ParentID   string        `json:"parent_id" swaggo:"false,父级ID"`
	ParentPath string        `json:"parent_path" swaggo:"false,父级路径"`
	Resources  MenuResources `json:"resources" swaggo:"false,资源列表"`
	Actions    MenuActions   `json:"actions" swaggo:"false,动作列表"`
	Children   *[]*MenuTree  `json:"children,omitempty" swaggo:"false,子级树"`
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

// ForEach 遍历数据项
func (a MenuTrees) ForEach(fn func(*MenuTree, int)) MenuTrees {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToTree 转换为树形结构
func (a MenuTrees) ToTree() []*MenuTree {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.RecordID] = item
	}

	var list []*MenuTree
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
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
	return list
}
