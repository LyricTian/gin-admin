package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
)

// NewTrans 创建事务管理实例
func NewTrans(db *gormplus.DB) *Trans {
	return &Trans{db}
}

// Trans 事务管理
type Trans struct {
	db *gormplus.DB
}

// Begin 开启事务
func (a *Trans) Begin(ctx context.Context) (interface{}, error) {
	result := a.db.Begin()
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return gormplus.Wrap(result), nil
}

// Commit 提交事务
func (a *Trans) Commit(ctx context.Context, trans interface{}) error {
	db, ok := trans.(*gormplus.DB)
	if !ok {
		return errors.New("unknow trans")
	}

	result := db.Commit()
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Rollback 回滚事务
func (a *Trans) Rollback(ctx context.Context, trans interface{}) error {
	db, ok := trans.(*gormplus.DB)
	if !ok {
		return errors.New("unknow trans")
	}

	result := db.Rollback()
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
