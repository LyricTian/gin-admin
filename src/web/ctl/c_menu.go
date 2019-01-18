package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Menu 菜单管理
type Menu struct {
	MenuBll *bll.Menu `inject:""`
}

// Query 查询数据
func (a *Menu) Query(ctx *context.Context) {
	switch ctx.Query("type") {
	case "page":
		a.QueryPage(ctx)
	case "tree":
		a.QueryTree(ctx)
	default:
		ctx.ResError(errors.NewBadRequestError("未知的查询类型"))
	}
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	params := schema.MenuPageQueryParam{
		Code:     ctx.Query("code"),
		Name:     ctx.Query("name"),
		ParentID: ctx.Query("parent_id"),
		Type:     util.S(ctx.Query("mtype")).Int(),
	}

	total, items, err := a.MenuBll.QueryPage(ctx.CContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResPage(total, items)
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(ctx *context.Context) {
	treeData, err := a.MenuBll.QueryTree(ctx.CContext())
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResList(treeData)
}

// Get 查询指定数据
func (a *Menu) Get(ctx *context.Context) {
	item, err := a.MenuBll.Get(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *Menu) Create(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	recordID, err := a.MenuBll.Create(ctx.CContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResSuccess(context.HTTPNewItem{RecordID: recordID})
}

// Update 更新数据
func (a *Menu) Update(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.MenuBll.Update(ctx.CContext(), ctx.Param("id"), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *Menu) Delete(ctx *context.Context) {
	err := a.MenuBll.Delete(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *Menu) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")
	if len(ids) == 0 {
		ctx.ResError(errors.NewBadRequestError("无效的请求数据"))
		return
	}

	err := a.MenuBll.Delete(ctx.CContext(), ids...)
	if err != nil {
		ctx.ResError(err)
		return
	}

	ctx.ResOK()
}
