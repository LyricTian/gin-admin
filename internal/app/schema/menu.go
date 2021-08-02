package schema

import (
	"strconv"
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v8/pkg/util/json"
)

// Menu 菜单对象
type Menu struct {
	ID         uint64      `json:"id"`                                     // 唯一标识
	Name       string      `json:"name" binding:"required"`                // 菜单名称
	Sequence   int         `json:"sequence"`                               // 排序值
	Icon       string      `json:"icon"`                                   // 菜单图标
	Router     string      `json:"router"`                                 // 访问路由
	ParentID   uint64      `json:"parent_id"`                              // 父级ID
	ParentPath string      `json:"parent_path"`                            // 父级路径
	IsShow     int         `json:"is_show" binding:"required,max=2,min=1"` // 是否显示(1:显示 2:隐藏)
	Status     int         `json:"status" binding:"required,max=2,min=1"`  // 状态(1:启用 2:禁用)
	Memo       string      `json:"memo"`                                   // 备注
	Creator    uint64      `json:"creator"`                                // 创建者
	CreatedAt  time.Time   `json:"created_at"`                             // 创建时间
	UpdatedAt  time.Time   `json:"updated_at"`                             // 更新时间
	Actions    MenuActions `json:"actions"`                                // 动作列表
}

func (a *Menu) String() string {
	return json.MarshalToString(a)
}

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	PaginationParam
	IDs              []uint64 `form:"-"`          // 唯一标识列表
	Name             string   `form:"-"`          // 菜单名称
	PrefixParentPath string   `form:"-"`          // 父级路径(前缀模糊查询)
	QueryValue       string   `form:"queryValue"` // 模糊查询
	ParentID         *uint64  `form:"parentID"`   // 父级内码
	IsShow           int      `form:"isShow"`     // 是否显示(1:显示 2:隐藏)
	Status           int      `form:"status"`     // 状态(1:启用 2:禁用)
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// MenuQueryResult 查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// Menus 菜单列表
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ToMap 转换为键值映射
func (a Menus) ToMap() map[uint64]*Menu {
	m := make(map[uint64]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// SplitParentIDs 拆分父级路径的唯一标识列表
func (a Menus) SplitParentIDs() []uint64 {
	idList := make([]uint64, 0, len(a))
	mIDList := make(map[uint64]struct{})

	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok || item.ParentPath == "" {
			continue
		}

		for _, pp := range strings.Split(item.ParentPath, "/") {
			ppID, _ := strconv.ParseUint(pp, 10, 64)
			if _, ok := mIDList[ppID]; ok {
				continue
			}
			idList = append(idList, ppID)
			mIDList[ppID] = struct{}{}
		}
	}

	return idList
}

// ToTree 转换为菜单树
func (a Menus) ToTree() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			ID:         item.ID,
			Name:       item.Name,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Sequence:   item.Sequence,
			IsShow:     item.IsShow,
			Status:     item.Status,
			Actions:    item.Actions,
		}
	}
	return list.ToTree()
}

// FillMenuAction 填充菜单动作列表
func (a Menus) FillMenuAction(mActions map[uint64]MenuActions) Menus {
	for _, item := range a {
		if v, ok := mActions[item.ID]; ok {
			item.Actions = v
		}
	}
	return a
}

// ----------------------------------------MenuTree--------------------------------------

// MenuTree 菜单树
type MenuTree struct {
	ID         uint64      `yaml:"-" json:"id"`                                  // 唯一标识
	Name       string      `yaml:"name" json:"name"`                             // 菜单名称
	Icon       string      `yaml:"icon" json:"icon"`                             // 菜单图标
	Router     string      `yaml:"router,omitempty" json:"router"`               // 访问路由
	ParentID   uint64      `yaml:"-" json:"parent_id"`                           // 父级ID
	ParentPath string      `yaml:"-" json:"parent_path"`                         // 父级路径
	Sequence   int         `yaml:"sequence" json:"sequence"`                     // 排序值
	IsShow     int         `yaml:"-" json:"is_show"`                             // 是否显示(1:显示 2:隐藏)
	Status     int         `yaml:"-" json:"status"`                              // 状态(1:启用 2:禁用)
	Actions    MenuActions `yaml:"actions,omitempty" json:"actions"`             // 动作列表
	Children   *MenuTrees  `yaml:"children,omitempty" json:"children,omitempty"` // 子级树
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

// ToTree 转换为树形结构
func (a MenuTrees) ToTree() MenuTrees {
	mi := make(map[uint64]*MenuTree)
	for _, item := range a {
		mi[item.ID] = item
	}

	var list MenuTrees
	for _, item := range a {
		if item.ParentID == 0 {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				children := MenuTrees{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

// ----------------------------------------MenuAction--------------------------------------

// MenuAction 菜单动作对象
type MenuAction struct {
	ID        uint64              `yaml:"-" json:"id"`                          // 唯一标识
	MenuID    uint64              `yaml:"-" binding:"required" json:"menu_id"`  // 菜单ID
	Code      string              `yaml:"code" binding:"required" json:"code"`  // 动作编号
	Name      string              `yaml:"name" binding:"required" json:"name"`  // 动作名称
	Resources MenuActionResources `yaml:"resources,omitempty" json:"resources"` // 资源列表
}

// MenuActionQueryParam 查询条件
type MenuActionQueryParam struct {
	PaginationParam
	MenuID uint64   // 菜单ID
	IDs    []uint64 // 唯一标识列表
}

// MenuActionQueryOptions 查询可选参数项
type MenuActionQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// MenuActionQueryResult 查询结果
type MenuActionQueryResult struct {
	Data       MenuActions
	PageResult *PaginationResult
}

// MenuActions 菜单动作管理列表
type MenuActions []*MenuAction

// ToMap 转换为map
func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

// FillResources 填充资源数据
func (a MenuActions) FillResources(mResources map[uint64]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.ID]
	}
}

// ToMenuIDMap 转换为菜单ID映射
func (a MenuActions) ToMenuIDMap() map[uint64]MenuActions {
	m := make(map[uint64]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}
	return m
}

// ----------------------------------------MenuActionResource--------------------------------------

// MenuActionResource 菜单动作关联资源对象
type MenuActionResource struct {
	ID       uint64 `yaml:"-" json:"id"`                             // 唯一标识
	ActionID uint64 `yaml:"-" json:"action_id"`                      // 菜单动作ID
	Method   string `yaml:"method" binding:"required" json:"method"` // 资源请求方式(支持正则)
	Path     string `yaml:"path" binding:"required" json:"path"`     // 资源请求路径（支持/:id匹配）
}

// MenuActionResourceQueryParam 查询条件
type MenuActionResourceQueryParam struct {
	PaginationParam
	MenuID  uint64   // 菜单ID
	MenuIDs []uint64 // 菜单ID列表
}

// MenuActionResourceQueryOptions 查询可选参数项
type MenuActionResourceQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// MenuActionResourceQueryResult 查询结果
type MenuActionResourceQueryResult struct {
	Data       MenuActionResources
	PageResult *PaginationResult
}

// MenuActionResources 菜单动作关联资源管理列表
type MenuActionResources []*MenuActionResource

// ToMap 转换为map
func (a MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range a {
		m[item.Method+item.Path] = item
	}
	return m
}

// ToActionIDMap 转换为动作ID映射
func (a MenuActionResources) ToActionIDMap() map[uint64]MenuActionResources {
	m := make(map[uint64]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}
	return m
}
