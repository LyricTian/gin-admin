package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// Demo 示例程序
type Demo struct {
	DemoModel model.IDemo `inject:"IDemo"`
	CommonBll *Common     `inject:""`
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx context.Context, params schema.DemoPageQueryParam, pageIndex, pageSize uint) (int, []*schema.Demo, error) {
	return a.DemoModel.QueryPage(ctx, params, pageIndex, pageSize)
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

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) (string, error) {
	exists, err := a.DemoModel.CheckCode(ctx, item.Code)
	if err != nil {
		return "", err
	} else if exists {
		return "", errors.NewBadRequestError("编号已经存在")
	}

	item.RecordID = util.MustUUID()
	err = a.DemoModel.Create(ctx, item)
	if err != nil {
		return "", err
	}
	return item.RecordID, nil
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item schema.Demo) error {
	oldItem, err := a.DemoModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Code != item.Code {
		exists, err := a.DemoModel.CheckCode(ctx, item.Code)
		if err != nil {
			return err
		} else if exists {
			return errors.NewBadRequestError("编号已经存在")
		}
	}

	return a.DemoModel.Update(ctx, recordID, item)
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordIDs ...string) error {
	return a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		for _, recordID := range recordIDs {
			err := a.DemoModel.Delete(ctx, recordID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	return a.DemoModel.UpdateStatus(ctx, recordID, status)
}
