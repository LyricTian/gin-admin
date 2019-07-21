package schema

import "time"

// Demo demo对象
type Demo struct {
	RecordID  string    `json:"record_id" swaggo:"false,记录ID"`
	Code      string    `json:"code" binding:"required" swaggo:"true,编号"`
	Name      string    `json:"name" binding:"required" swaggo:"true,名称"`
	Memo      string    `json:"memo" swaggo:"false,备注"`
	Status    int       `json:"status" binding:"required,max=2,min=1" swaggo:"true,状态(1:启用 2:停用)"`
	Creator   string    `json:"creator" swaggo:"false,创建者"`
	CreatedAt time.Time `json:"created_at" swaggo:"false,创建时间"`
}

// DemoQueryParam 查询条件
type DemoQueryParam struct {
	Code     string // 编号
	Status   int    // 状态(1:启用 2:停用)
	LikeCode string // 编号(模糊查询)
	LikeName string // 名称(模糊查询)
}

// DemoQueryOptions demo对象查询可选参数项
type DemoQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// DemoQueryResult demo对象查询结果
type DemoQueryResult struct {
	Data       []*Demo
	PageResult *PaginationResult
}
