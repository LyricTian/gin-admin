package schema

import "time"

// Resource 资源对象
type Resource struct {
	RecordID  string    `json:"record_id" swaggo:"false,记录ID"`
	Code      string    `json:"code" binding:"required" swaggo:"true,资源编号"`
	Name      string    `json:"name" binding:"required" swaggo:"true,资源名称"`
	Path      string    `json:"path" swaggo:"true,访问路径"`
	Method    string    `json:"method" swaggo:"true,资源请求方式"`
	Creator   string    `json:"creator" swaggo:"false,创建者"`
	CreatedAt time.Time `json:"created_at" swaggo:"false,创建时间"`
}

// ResourceQueryParam 查询条件
type ResourceQueryParam struct {
	Name string // 资源名称（模糊查询）
	Path string // 访问路径（前缀模糊查询）
}

// ResourceQueryOptions 查询可选参数项
type ResourceQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// ResourceQueryResult 查询结果
type ResourceQueryResult struct {
	Data       Resources
	PageResult *PaginationResult
}

// Resources 资源对象列表
type Resources []*Resource
