package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Demo 示例程序
type Demo struct {
	DemoBll *bll.Demo `inject:""`
}

// Query 查询数据
func (a *Demo) Query(ctx *context.Context) {
	switch ctx.Query("type") {
	case "page":
		a.QueryPage(ctx)
	default:
		ctx.ResBadRequest(nil)
	}
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	var params schema.DemoQueryParam

	params.Code = ctx.Query("code")
	params.Name = ctx.Query("name")

	total, items, err := a.DemoBll.QueryPage(ctx.NewContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResPage(total, items)
}

// Get 查询指定数据
func (a *Demo) Get(ctx *context.Context) {
	item, err := a.DemoBll.Get(ctx.NewContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *Demo) Create(ctx *context.Context) {
	var item schema.Demo
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	item.Creator = ctx.GetUserID()
	err := a.DemoBll.Create(ctx.NewContext(), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	newItem, err := a.DemoBll.Get(ctx.NewContext(), item.RecordID)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResSuccess(newItem)
}

// Update 更新数据
func (a *Demo) Update(ctx *context.Context) {
	var item schema.Demo
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.DemoBll.Update(ctx.NewContext(), ctx.Param("id"), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResOK()
}

// Delete 删除数据
func (a *Demo) Delete(ctx *context.Context) {
	err := a.DemoBll.Delete(ctx.NewContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *Demo) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")

	for _, id := range ids {
		err := a.DemoBll.Delete(ctx.NewContext(), id)
		if err != nil {
			ctx.ResInternalServerError(err)
			return
		}
	}

	ctx.ResOK()
}
