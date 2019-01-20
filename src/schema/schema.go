package schema

// PaginationParam 分页查询条件
type PaginationParam struct {
	PageIndex uint // 页索引
	PageSize  uint // 页大小
}

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total int // 总数据条数
}
