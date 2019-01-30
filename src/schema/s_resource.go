package schema

// Resource 资源对象
type Resource struct {
	RecordID string `json:"record_id" swaggo:"false,记录ID"`
	Code     string `json:"code" binding:"required" swaggo:"true,资源编号"`
	Name     string `json:"name" binding:"required" swaggo:"true,资源名称"`
	Path     string `json:"path" swaggo:"true,访问路径"`
	Method   string `json:"method" swaggo:"true,资源请求方式"`
}

// ResourceQueryParam 资源对象查询条件
type ResourceQueryParam struct {
	Name string // 资源名称（模糊查询）
	Path string // 访问路径（前缀模糊查询）
}

// ResourceQueryOptions 资源对象查询可选参数项
type ResourceQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// ResourceQueryResult 资源对象查询结果
type ResourceQueryResult struct {
	Data       []*Resource
	PageResult *PaginationResult
}
