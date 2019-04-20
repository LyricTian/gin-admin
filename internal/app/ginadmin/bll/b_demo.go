package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/util"
)

// NewDemo 创建demo
func NewDemo(m *model.Common) *Demo {
	return &Demo{
		DemoModel: m.Demo,
	}
}

// Demo 示例程序
type Demo struct {
	DemoModel model.IDemo
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx context.Context, params schema.DemoQueryParam, pp *schema.PaginationParam) ([]*schema.Demo, *schema.PaginationResult, error) {
	result, err := a.DemoModel.Query(ctx, params, schema.DemoQueryOptions{PageParam: pp})
	if err != nil {
		return nil, nil, err
	}
	return result.Data, result.PageResult, nil
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string) (*schema.Demo, error) {
	item, err := a.DemoModel.Get(ctx, recordID)
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
		return errors.NewBadRequestError("编号已经存在")
	}
	return nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) (*schema.Demo, error) {
	err := a.checkCode(ctx, item.Code)
	if err != nil {
		return nil, err
	}

	item.RecordID = util.MustUUID()
	item.Creator = GetUserID(ctx)
	err = a.DemoModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.Get(ctx, item.RecordID)
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
	return a.Get(ctx, recordID)
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	err := a.DemoModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	return a.DemoModel.UpdateStatus(ctx, recordID, status)
}
