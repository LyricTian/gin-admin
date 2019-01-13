package schema

import "time"

// Demo 示例程序
type Demo struct {
	RecordID string `json:"record_id"`                 // 记录内码
	Code     string `json:"code" binding:"required"`   // 编号
	Name     string `json:"name" binding:"required"`   // 名称
	Memo     string `json:"memo"`                      // 备注
	Status   int    `json:"status" binding:"required"` // 状态(1:启用 2:停用)
	Creator  string `json:"creator"`                   // 创建者
}

// DemoQueryParam 示例查询条件
type DemoQueryParam struct {
	Code   string // 编号
	Name   string // 名称
	Status int    // 状态(1:启用 2:停用)
}

// DemoQueryResult 示例查询结果
type DemoQueryResult struct {
	RecordID string    `json:"record_id"`  // 记录内码
	Code     string    `json:"code"`       // 编号
	Name     string    `json:"name"`       // 名称
	Memo     string    `json:"memo"`       // 备注
	Status   int       `json:"status"`     // 状态(1:启用 2:停用)
	Created  time.Time `json:"created_at"` // 创建时间
}
