package api

import (
	"gin-admin/src/bll"
	"gin-admin/src/context"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"strings"
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
	default:
		ctx.ResBadRequest(nil)
	}
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	params := schema.MenuQueryParam{
		Name:     ctx.Query("name"),
		ParentID: ctx.Query("parent_id"),
		Status:   util.S(ctx.Query("status")).Int(),
		Type:     util.S(ctx.Query("type")).Int(),
	}

	total, items, err := a.MenuBll.QueryPage(params, pageIndex, pageSize)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResPage(total, items)
}

// Get 查询指定数据
func (a *Menu) Get(ctx *context.Context) {
	item, err := a.MenuBll.Get(ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *Menu) Create(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	item.Creator = ctx.GetUserID()
	err := a.MenuBll.Create(&item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Update 更新数据
func (a *Menu) Update(ctx *context.Context) {
	var item schema.Menu
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.MenuBll.Update(ctx.Param("id"), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *Menu) Delete(ctx *context.Context) {
	err := a.MenuBll.Delete(ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *Menu) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")

	for _, id := range ids {
		err := a.MenuBll.Delete(id)
		if err != nil {
			ctx.ResInternalServerError(err)
			return
		}
	}

	ctx.ResOK()
}
