package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Role 角色管理
type Role struct {
	RoleBll *bll.Role `inject:""`
}

// Query 查询数据
func (a *Role) Query(ctx *context.Context) {
	switch ctx.Query("type") {
	case "page":
		a.QueryPage(ctx)
	case "select":
		a.QuerySelect(ctx)
	default:
		ctx.ResBadRequest(nil)
	}
}

// QueryPage 查询分页数据
func (a *Role) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	var params schema.RoleQueryParam

	params.Name = ctx.Query("name")
	params.Status = util.S(ctx.Query("status")).Int()

	total, items, err := a.RoleBll.QueryPage(ctx.NewContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResPage(total, items)
}

// QuerySelect 查询分页数据
func (a *Role) QuerySelect(ctx *context.Context) {
	var params schema.RoleSelectQueryParam

	params.Name = ctx.Query("name")
	params.Status = util.S(ctx.Query("status")).Int()

	items, err := a.RoleBll.QuerySelect(ctx.NewContext(), params)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResList(items)
}

// Get 查询指定数据
func (a *Role) Get(ctx *context.Context) {
	item, err := a.RoleBll.Get(ctx.NewContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *Role) Create(ctx *context.Context) {
	var item schema.Role
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	item.Creator = ctx.GetUserID()
	err := a.RoleBll.Create(ctx.NewContext(), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	newItem, err := a.RoleBll.Get(ctx.NewContext(), item.RecordID)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResSuccess(newItem)
}

// Update 更新数据
func (a *Role) Update(ctx *context.Context) {
	var item schema.Role
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.RoleBll.Update(ctx.NewContext(), ctx.Param("id"), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *Role) Delete(ctx *context.Context) {
	err := a.RoleBll.Delete(ctx.NewContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *Role) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")

	for _, id := range ids {
		err := a.RoleBll.Delete(ctx.NewContext(), id)
		if err != nil {
			ctx.ResInternalServerError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
func (a *Role) Enable(ctx *context.Context) {
	err := a.RoleBll.UpdateStatus(ctx.NewContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
func (a *Role) Disable(ctx *context.Context) {
	err := a.RoleBll.UpdateStatus(ctx.NewContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}
