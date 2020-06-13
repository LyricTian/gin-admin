package api

import (
	"strings"

	"github.com/LyricTian/gin-admin/v6/internal/app/bll"
	"github.com/LyricTian/gin-admin/v6/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/v6/internal/app/schema"
	"github.com/LyricTian/gin-admin/v6/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User 用户管理
type User struct {
	UserBll bll.IUser
}

// Query 查询数据
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	params.Pagination = true
	result, err := a.UserBll.QueryShow(ctx, params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserBll.Get(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item.CleanSecure())
}

// Create 创建数据
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	} else if item.Password == "" {
		ginplus.ResError(c, errors.New400Response("密码不能为空"))
		return
	}

	item.Creator = ginplus.GetUserID(c)
	result, err := a.UserBll.Create(ctx, item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// Update 更新数据
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.UserBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Delete 删除数据
func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBll.Delete(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable 启用数据
func (a *User) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBll.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable 禁用数据
func (a *User) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBll.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
