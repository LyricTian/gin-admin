package schema

// HTTPStatusText 定义HTTP状态文本
type HTTPStatusText string

func (t HTTPStatusText) String() string {
	return string(t)
}

// 定义HTTP状态文本常量
const (
	OKStatusText HTTPStatusText = "OK"
)

// HTTPError HTTP响应错误
type HTTPError struct {
	Error HTTPErrorItem `json:"error"` // 错误项
}

// HTTPErrorItem HTTP响应错误项
type HTTPErrorItem struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// HTTPStatus HTTP响应状态
type HTTPStatus struct {
	Status string `json:"status"` // 状态(OK)
}

// HTTPList HTTP响应列表数据
type HTTPList struct {
	List       interface{}     `json:"list"`
	Pagination *HTTPPagination `json:"pagination,omitempty"`
}

// HTTPPagination HTTP分页数据
type HTTPPagination struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
}

// NewPaginationParam 创建分页查询条件实例
func NewPaginationParam(pageSize int, pageIndex ...int) *PaginationParam {
	item := &PaginationParam{
		PageSize: pageSize,
	}

	if len(pageIndex) > 0 {
		item.PageIndex = pageIndex[0]
	}

	return item
}

// PaginationParam 分页查询条件
type PaginationParam struct {
	PageIndex int // 页索引
	PageSize  int // 页大小
}

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total int // 总数据条数
}

// NewHTTPRecordID 创建HTTP记录ID实例
func NewHTTPRecordID(recordID string) *HTTPRecordID {
	return &HTTPRecordID{
		RecordID: recordID,
	}
}

// HTTPRecordID HTTP记录ID
type HTTPRecordID struct {
	RecordID string `json:"record_id"`
}

// ----------------------------------------OrderField--------------------------------------

// OrderDirection 排序方向
type OrderDirection int

const (
	// OrderByASC 升序排序
	OrderByASC OrderDirection = 1
	// OrderByDESC 降序排序
	OrderByDESC OrderDirection = 2
)

// NewOrderFields 创建排序字段(默认升序排序)，可指定不同key的排序规则
// keys 需要排序的key
// directions 排序规则，按照key的索引指定，索引默认从0开始
func NewOrderFields(keys []string, directions ...map[int]OrderDirection) []*OrderField {
	m := make(map[int]OrderDirection)
	if len(directions) > 0 {
		m = directions[0]
	}

	fields := make([]*OrderField, len(keys))
	for i, key := range keys {
		d := OrderByASC
		if v, ok := m[i]; ok {
			d = v
		}

		fields[i] = NewOrderField(key, d)
	}

	return fields
}

// NewOrderField 创建排序字段
func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

// OrderField 排序字段
type OrderField struct {
	Key       string         // 字段名(字段名约束为小写蛇形)
	Direction OrderDirection // 排序方向
}
