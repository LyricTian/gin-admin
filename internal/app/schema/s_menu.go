package schema

import (
	"strings"
	"time"
)

// Menu 菜单对象
type Menu struct {
	RecordID   string      `json:"record_id"`                             // 记录ID
	Name       string      `json:"name" binding:"required"`               // 菜单名称
	Sequence   int         `json:"sequence"`                              // 排序值
	Icon       string      `json:"icon"`                                  // 菜单图标
	Router     string      `json:"router"`                                // 访问路由
	ParentID   string      `json:"parent_id"`                             // 父级ID
	ParentPath string      `json:"parent_path"`                           // 父级路径
	Status     int         `json:"status" binding:"required,max=2,min=1"` // 状态(1:正常 2:隐藏)
	Memo       string      `json:"memo"`                                  // 备注
	Creator    string      `json:"creator"`                               // 创建者
	CreatedAt  time.Time   `json:"created_at"`                            // 创建时间
	UpdatedAt  time.Time   `json:"updated_at"`                            // 更新时间
	Actions    MenuActions `json:"actions"`                               // 动作列表
}

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	RecordIDs        []string // 记录ID列表
	LikeName         string   // 菜单名称(模糊查询)
	Name             string   // 菜单名称
	ParentID         *string  // 父级内码
	PrefixParentPath string   // 父级路径(前缀模糊查询)
	Status           int      // 状态(1:正常 2:隐藏)
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
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
	recordIDs := make([]string, 0, len(a))
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

// ToTree 转换为菜单树
func (a Menus) ToTree() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			RecordID:   item.RecordID,
			Name:       item.Name,
			Sequence:   item.Sequence,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Status:     item.Status,
			Memo:       item.Memo,
			Actions:    item.Actions,
		}
	}
	return list.ToTree()
}

func (a Menus) fillLeafNodeID(tree *MenuTrees, leafNodeIDs *[]string) {
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
	tree := a.ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

// MenuTree 菜单树
type MenuTree struct {
	RecordID   string      `json:"record_id"`          // 记录ID
	Name       string      `json:"name"`               // 菜单名称
	Sequence   int         `json:"sequence"`           // 排序值
	Icon       string      `json:"icon"`               // 菜单图标
	Router     string      `json:"router"`             // 访问路由
	ParentID   string      `json:"parent_id"`          // 父级ID
	ParentPath string      `json:"parent_path"`        // 父级路径
	Status     int         `json:"status"`             // 状态(1:正常 2:隐藏)
	Memo       string      `json:"memo"`               // 备注
	Actions    MenuActions `json:"actions"`            // 动作列表
	Children   *MenuTrees  `json:"children,omitempty"` // 子级树
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

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
				var children MenuTrees
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}
