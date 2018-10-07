package api

import (
	"gin-admin/src/bll"
	"gin-admin/src/context"
	"gin-admin/src/schema"
	"gin-admin/src/util"
)

// Demo 示例
type Demo struct {
	DemoBll *bll.Demo `inject:""`
}

// QueryList /
func (a *Demo) QueryList(ctx *context.Context) {
	items, err := a.DemoBll.Query()
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResList(items)
}

// Get /:id
func (a *Demo) Get(ctx *context.Context) {
	id := util.S(ctx.Param("id")).Int64()
	item, err := a.DemoBll.Get(id)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建实例数据
func (a *Demo) Create(ctx *context.Context) {
	var reqData schema.Demo
	if err := ctx.ParseJSON(&reqData); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.DemoBll.Create(&reqData)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Update 更新实例数据
func (a *Demo) Update(ctx *context.Context) {
	var reqData schema.Demo
	if err := ctx.ParseJSON(&reqData); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.DemoBll.Update(util.S(ctx.Param("id")).Int64(), &reqData)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除实例数据
func (a *Demo) Delete(ctx *context.Context) {
	err := a.DemoBll.Delete(util.S(ctx.Param("id")).Int64())
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}
