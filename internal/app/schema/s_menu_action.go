package schema

// MenuAction 菜单动作管理
type MenuAction struct {
	RecordID  string              `json:"record_id"`                  // 记录ID
	MenuID    string              `json:"menu_id" binding:"required"` // 菜单ID
	Code      string              `json:"code" binding:"required"`    // 动作编号
	Name      string              `json:"name" binding:"required"`    // 动作名称
	Resources MenuActionResources `json:"resources"`                  // 资源列表
}

// MenuActionQueryParam 查询条件
type MenuActionQueryParam struct {
	MenuID string // 菜单ID
}

// MenuActionQueryOptions 查询可选参数项
type MenuActionQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
}

// MenuActionQueryResult 查询结果
type MenuActionQueryResult struct {
	Data       MenuActions
	PageResult *PaginationResult
}

// MenuActions 菜单动作管理列表
type MenuActions []*MenuAction

// FillResources 填充资源数据
func (a MenuActions) FillResources(mResources map[string]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.RecordID]
	}
}

// GetByRecordID 根据记录ID获取数据项
func (a MenuActions) GetByRecordID(recordID string) *MenuAction {
	for _, item := range a {
		if item.RecordID == recordID {
			return item
		}
	}
	return nil
}
