package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Demo 示例程序
// @Name Demo
// @Description 示例程序
type Demo struct {
	DemoBll *bll.Demo `inject:""`
}

// Query 查询数据
// @Title 查询数据
// @Description 查询示例数据(支持分页查询)
// @Consumes json
// @Produces json
// @Param access-token header string false "访问令牌"
// @Param type query string true "查询类型(分页查询：page)"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Param code query string false "编号"
// @Param name query string false "名称"
// @Param status query int false "状态"
// @Success 200 []schema.Demo "分页查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小}}"
// @Failure 400 string "{error:{code:0,message:未知的查询类型}}"
// @Failure 500 string "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/demos
func (a *Demo) Query(ctx *context.Context) {
	switch ctx.Query("type") {
	case "page":
		a.QueryPage(ctx)
	default:
		ctx.ResError(errors.NewBadRequestError("未知的查询类型"))
	}
}

// QueryPage 查询分页数据
func (a *Demo) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	var params schema.DemoPageQueryParam
	params.Code = ctx.Query("code")
	params.Name = ctx.Query("name")
	params.Status = util.S(ctx.Query("status")).Int()

	total, items, err := a.DemoBll.QueryPage(ctx.CContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResPage(total, items)
}

// Get 查询指定数据
func (a *Demo) Get(ctx *context.Context) {
	item, err := a.DemoBll.Get(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *Demo) Create(ctx *context.Context) {
	var item schema.Demo
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	recordID, err := a.DemoBll.Create(ctx.CContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResSuccess(context.HTTPNewItem{RecordID: recordID})
}

// Update 更新数据
func (a *Demo) Update(ctx *context.Context) {
	var item schema.Demo
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.DemoBll.Update(ctx.CContext(), ctx.Param("id"), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *Demo) Delete(ctx *context.Context) {
	err := a.DemoBll.Delete(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *Demo) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")

	for _, id := range ids {
		err := a.DemoBll.Delete(ctx.CContext(), id)
		if err != nil {
			ctx.ResError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
func (a *Demo) Enable(ctx *context.Context) {
	err := a.DemoBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
func (a *Demo) Disable(ctx *context.Context) {
	err := a.DemoBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}
