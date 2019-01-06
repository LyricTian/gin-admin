package ctl

import (
	"strings"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// User 用户管理
type User struct {
	UserBll *bll.User `inject:""`
}

// Query 查询数据
func (a *User) Query(ctx *context.Context) {
	switch ctx.Query("type") {
	case "page":
		a.QueryPage(ctx)
	default:
		ctx.ResBadRequest(nil)
	}
}

// QueryPage 查询分页数据
func (a *User) QueryPage(ctx *context.Context) {
	pageIndex, pageSize := ctx.GetPageIndex(), ctx.GetPageSize()

	var params schema.UserQueryParam

	params.UserName = ctx.Query("user_name")
	params.RealName = ctx.Query("real_name")
	params.RoleID = ctx.Query("role_id")
	params.Status = util.S(ctx.Query("status")).Int()

	total, items, err := a.UserBll.QueryPage(ctx.CContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResPage(total, items)
}

// Get 查询指定数据
func (a *User) Get(ctx *context.Context) {
	item, err := a.UserBll.Get(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResSuccess(item)
}

// Create 创建数据
func (a *User) Create(ctx *context.Context) {
	var item schema.User
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	item.Creator = ctx.GetUserID()
	err := a.UserBll.Create(ctx.CContext(), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	newItem, err := a.UserBll.Get(ctx.CContext(), item.RecordID)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResSuccess(newItem)
}

// Update 更新数据
func (a *User) Update(ctx *context.Context) {
	var item schema.User
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.UserBll.Update(ctx.CContext(), ctx.Param("id"), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *User) Delete(ctx *context.Context) {
	err := a.UserBll.Delete(ctx.CContext(), ctx.Param("id"))
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// DeleteMany 删除多条数据
func (a *User) DeleteMany(ctx *context.Context) {
	ids := strings.Split(ctx.Query("batch"), ",")

	for _, id := range ids {
		err := a.UserBll.Delete(ctx.CContext(), id)
		if err != nil {
			ctx.ResInternalServerError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
func (a *User) Enable(ctx *context.Context) {
	err := a.UserBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
func (a *User) Disable(ctx *context.Context) {
	err := a.UserBll.UpdateStatus(ctx.CContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}
