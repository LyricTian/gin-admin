package bll

import (
	"context"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"time"
)

// Demo 示例程序
type Demo struct {
	DemoModel model.IDemo `inject:"IDemo"`
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx context.Context, params schema.DemoQueryParam, pageIndex, pageSize uint) (int64, []*schema.DemoQueryResult, error) {
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
func (a *Demo) Create(ctx context.Context, item *schema.Demo) error {
	item.ID = 0
	item.RecordID = util.MustUUID()
	item.Created = time.Now().Unix()
	item.Deleted = 0
	return a.DemoModel.Create(ctx, item)
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item *schema.Demo) error {
	exists, err := a.DemoModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	info := util.StructToMap(item)
	delete(info, "id")
	delete(info, "record_id")
	delete(info, "creator")
	delete(info, "created")
	delete(info, "deleted")

	return a.DemoModel.Update(ctx, recordID, info)
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	exists, err := a.DemoModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	return a.DemoModel.Delete(ctx, recordID)
}
