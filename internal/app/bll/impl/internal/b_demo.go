package internal

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
)

// NewDemo 创建demo
func NewDemo(mDemo model.IDemo) *Demo {
	return &Demo{
		DemoModel: mDemo,
	}
}

// Demo 示例程序
type Demo struct {
	DemoModel model.IDemo
}

// Query 查询数据
func (a *Demo) Query(ctx context.Context, params schema.DemoQueryParam, opts ...schema.DemoQueryOptions) (*schema.DemoQueryResult, error) {
	return a.DemoModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string, opts ...schema.DemoQueryOptions) (*schema.Demo, error) {
	item, err := a.DemoModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Demo) checkCode(ctx context.Context, code string) error {
	result, err := a.DemoModel.Query(ctx, schema.DemoQueryParam{
		Code: code,
	}, schema.DemoQueryOptions{
		PageParam: &schema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.ErrResourceExists
	}
	return nil
}

func (a *Demo) getUpdate(ctx context.Context, recordID string) (*schema.Demo, error) {
	return a.Get(ctx, recordID)
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) (*schema.Demo, error) {
	err := a.checkCode(ctx, item.Code)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.MustUUID()
	err = a.DemoModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.RecordID)
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item schema.Demo) (*schema.Demo, error) {
	oldItem, err := a.DemoModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Code != item.Code {
		err := a.checkCode(ctx, item.Code)
		if err != nil {
			return nil, err
		}
	}

	err = a.DemoModel.Update(ctx, recordID, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, recordID)
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	oldItem, err := a.DemoModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.DemoModel.Delete(ctx, recordID)
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	oldItem, err := a.DemoModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.DemoModel.UpdateStatus(ctx, recordID, status)
}
