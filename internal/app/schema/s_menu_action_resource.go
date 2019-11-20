package schema

// MenuActionResource 菜单动作关联资源管理
type MenuActionResource struct {
	RecordID string `json:"record_id"`                    // 记录ID
	ActionID string `json:"action_id" binding:"required"` // 菜单动作ID
	Method   string `json:"method" binding:"required"`    // 资源请求方式(支持正则)
	Path     string `json:"path" binding:"required"`      // 资源请求路径（支持/:id匹配）
}

// MenuActionResourceQueryParam 查询条件
type MenuActionResourceQueryParam struct {
	MenuID string // 菜单ID
}

// MenuActionResourceQueryOptions 查询可选参数项
type MenuActionResourceQueryOptions struct {
	PageParam   *PaginationParam // 分页参数
	OrderFields []*OrderField    // 排序字段
}

// MenuActionResourceQueryResult 查询结果
type MenuActionResourceQueryResult struct {
	Data       MenuActionResources
	PageResult *PaginationResult
}

// MenuActionResources 菜单动作关联资源管理列表
type MenuActionResources []*MenuActionResource

// ToActionIDMap 转换为动作ID映射
func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)

	for _, item := range a {
		if v, ok := m[item.ActionID]; ok {
			v = append(v, item)
			m[item.ActionID] = v
			continue
		}
		m[item.ActionID] = MenuActionResources{item}
	}

	return m
}

// GetByRecordID 根据记录ID获取数据项
func (a MenuActionResources) GetByRecordID(recordID string) *MenuActionResource {
	for _, item := range a {
		if item.RecordID == recordID {
			return item
		}
	}
	return nil
}
