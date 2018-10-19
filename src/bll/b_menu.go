package bll

import (
	"context"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"time"
)

// Menu 菜单管理
type Menu struct {
	MenuModel model.IMenu `inject:"IMenu"`
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx context.Context, params schema.MenuQueryParam, pageIndex, pageSize uint) (int64, []*schema.MenuQueryResult, error) {
	return a.MenuModel.QueryPage(ctx, params, pageIndex, pageSize)
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string) (*schema.Menu, error) {
	return a.MenuModel.Get(ctx, recordID)
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item *schema.Menu) error {
	item.ID = 0
	item.RecordID = util.UUIDString()
	item.Created = time.Now().Unix()
	item.Deleted = 0
	return a.MenuModel.Create(ctx, item)
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item *schema.Menu) error {
	info := util.StructToMap(item)
	delete(info, "id")
	delete(info, "record_id")
	delete(info, "creator")
	delete(info, "created")
	delete(info, "deleted")

	return a.MenuModel.Update(ctx, recordID, info)
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	return a.MenuModel.Delete(ctx, recordID)
}
