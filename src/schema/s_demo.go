package schema

// Demo demo
type Demo struct {
	RecordID string `json:"record_id" swaggo:"false,记录ID"`              // 记录ID
	Code     string `json:"code" binding:"required" swaggo:"true,编号"`   // 编号
	Name     string `json:"name" binding:"required" swaggo:"true,名称"`   // 名称
	Memo     string `json:"memo" swaggo:"false,备注"`                     // 备注
	Status   int    `json:"status" binding:"required" swaggo:"true,状态"` // 状态(1:启用 2:停用)
}

// DemoPageQueryParam 分页查询条件
type DemoPageQueryParam struct {
	Code   string // 编号
	Name   string // 名称
	Status int    // 状态(1:启用 2:停用)
}
