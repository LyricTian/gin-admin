package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/google/wire"
)

var _ bll.IDemo = (*Demo)(nil)

// DemoSet 注入Demo
var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"), wire.Bind(new(bll.IDemo), new(*Demo)))

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
		PaginationParam: schema.PaginationParam{
			OnlyCount: true,
		},
		Code: code,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("编号已经存在")
	}

	return nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) (*schema.RecordIDResult, error) {
	err := a.checkCode(ctx, item.Code)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.NewRecordID()
	err = a.DemoModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return schema.NewRecordIDResult(item.RecordID), nil
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item schema.Demo) error {
	oldItem, err := a.DemoModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Code != item.Code {
		if err := a.checkCode(ctx, item.Code); err != nil {
			return err
		}
	}
	item.RecordID = oldItem.RecordID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	return a.DemoModel.Update(ctx, recordID, item)
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
