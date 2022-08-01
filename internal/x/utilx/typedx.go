package utilx

import "github.com/LyricTian/gin-admin/v9/pkg/errors"

const (
	RequestBodyKey  = "req-body"
	ResponseBodyKey = "res-body"
)

type OkResult struct {
	Ok bool `json:"ok"`
}

type ErrorResult struct {
	Error *errors.Error `json:"error"`
}

type ListResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

type PaginationParam struct {
	Pagination bool `form:"-"`
	OnlyCount  bool `form:"-"`
	Current    int  `form:"current,default=1"`
	PageSize   int  `form:"pageSize,default=10" binding:"max=100"`
}

func (a PaginationParam) GetCurrent() int {
	return a.Current
}

func (a PaginationParam) GetPageSize() int {
	pageSize := a.PageSize
	if a.PageSize <= 0 {
		pageSize = 100
	}
	return pageSize
}

type QueryOptions struct {
	SelectFields []string
	OmitFields   []string
	OrderFields  OrderByParams
}

type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)

type OrderByParam struct {
	Field     string
	Direction Direction
}

type OrderByParams []OrderByParam

func (a OrderByParams) ToSQL() string {
	if len(a) == 0 {
		return ""
	}

	var sql string
	for _, v := range a {
		sql += v.Field + " " + string(v.Direction) + ","
	}
	return sql[:len(sql)-1]
}
