package api

import (
	"gin-admin/src/bll"
	"gin-admin/src/context"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"strings"
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

	total, items, err := a.UserBll.QueryPage(ctx.NewContext(), params, pageIndex, pageSize)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}

	ctx.ResPage(total, items)
}

// Get 查询指定数据
func (a *User) Get(ctx *context.Context) {
	item, err := a.UserBll.Get(ctx.NewContext(), ctx.Param("id"))
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
	err := a.UserBll.Create(ctx.NewContext(), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Update 更新数据
func (a *User) Update(ctx *context.Context) {
	var item schema.User
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	err := a.UserBll.Update(ctx.NewContext(), ctx.Param("id"), &item)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Delete 删除数据
func (a *User) Delete(ctx *context.Context) {
	err := a.UserBll.Delete(ctx.NewContext(), ctx.Param("id"))
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
		err := a.UserBll.Delete(ctx.NewContext(), id)
		if err != nil {
			ctx.ResInternalServerError(err)
			return
		}
	}

	ctx.ResOK()
}

// Enable 启用数据
func (a *User) Enable(ctx *context.Context) {
	err := a.UserBll.UpdateStatus(ctx.NewContext(), ctx.Param("id"), 1)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}

// Disable 禁用数据
func (a *User) Disable(ctx *context.Context) {
	err := a.UserBll.UpdateStatus(ctx.NewContext(), ctx.Param("id"), 2)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResOK()
}
