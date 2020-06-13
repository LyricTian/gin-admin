package model

import (
	"context"
)

// ITrans 事务管理接口
type ITrans interface {
	// 执行事务
	Exec(ctx context.Context, fn func(context.Context) error) error
}
