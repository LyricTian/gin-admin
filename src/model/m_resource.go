package model

import (
	"context"

	"github.com/LyricTian/gin-admin/src/schema"
)

// IResource 资源存储接口
type IResource interface {
	// 查询数据
	Query(ctx context.Context, params schema.ResourceQueryParam, opts ...schema.ResourceQueryOptions) (schema.ResourceQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string) (*schema.Resource, error)
	// 检查访问路径和请求方法是否存在
	CheckPathAndMethod(ctx context.Context, path, method string) (bool, error)
	// 创建数据
	Create(ctx context.Context, item schema.Resource) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Resource) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
}
