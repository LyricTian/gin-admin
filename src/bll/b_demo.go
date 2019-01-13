package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// Demo 示例程序
type Demo struct {
	DemoModel  model.IDemo  `inject:"IDemo"`
	TransModel model.ITrans `inject:"ITrans"`
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx context.Context, params schema.DemoQueryParam, pageIndex, pageSize uint) (int, []schema.DemoQueryResult, error) {
	return a.DemoModel.QueryPage(ctx, params, pageIndex, pageSize)
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string) (*schema.Demo, error) {
	item, err := a.DemoModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrNotFound
	}

	return item, nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) (string, error) {
	exists, err := a.DemoModel.CheckCode(ctx, item.Code)
	if err != nil {
		return "", err
	} else if exists {
		return "", util.NewBadRequestError("编号已经存在")
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
		return util.ErrNotFound
	} else if oldItem.Code != item.Code {
		exists, err := a.DemoModel.CheckCode(ctx, item.Code)
		if err != nil {
			return err
		} else if exists {
			return util.NewBadRequestError("编号已经存在")
		}
	}

	return a.DemoModel.Update(ctx, recordID, item)
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordIDs ...string) error {
	trans, err := a.TransModel.Begin(ctx)
	if err != nil {
		return err
	}

	for _, recordID := range recordIDs {
		exists, err := a.DemoModel.Check(ctx, recordID)
		if err != nil {
			a.TransModel.Rollback(ctx, trans)
			return err
		} else if !exists {
			a.TransModel.Rollback(ctx, trans)
			return util.ErrNotFound
		}

		err = a.DemoModel.Delete(ctx, trans, recordID)
		if err != nil {
			a.TransModel.Rollback(ctx, trans)
			return err
		}
	}

	err = a.TransModel.Commit(ctx, trans)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	exists, err := a.DemoModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	return a.DemoModel.UpdateStatus(ctx, recordID, status)
}
