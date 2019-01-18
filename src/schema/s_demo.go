package schema

// Demo demo
type Demo struct {
	RecordID string `json:"record_id"`                 // 记录内码
	Code     string `json:"code" binding:"required"`   // 编号
	Name     string `json:"name" binding:"required"`   // 名称
	Memo     string `json:"memo"`                      // 备注
	Status   int    `json:"status" binding:"required"` // 状态(1:启用 2:停用)
}

// DemoPageQueryParam 分页查询条件
type DemoPageQueryParam struct {
	Code   string // 编号
	Name   string // 名称
	Status int    // 状态(1:启用 2:停用)
}
