package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/LyricTian/gin-admin/v8/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/internal/app/service"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
	"github.com/LyricTian/gin-admin/v8/pkg/util/conv"
)

var UserSet = wire.NewSet(wire.Struct(new(UserAPI), "*"))

// UserAPI 用户管理
type UserAPI struct {
	UserSrv *service.UserSrv
}

// Query 查询数据
func (a *UserAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}
	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = conv.ParseStringSliceToUint64(strings.Split(v, ","))
	}

	params.Pagination = true
	result, err := a.UserSrv.QueryShow(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
func (a *UserAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserSrv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item.CleanSecure())
}

// Create 创建数据
func (a *UserAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	} else if item.Password == "" {
		ginx.ResError(c, errors.New400Response("密码不能为空"))
		return
	}

	item.Creator = contextx.FromUserID(ctx)
	result, err := a.UserSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
func (a *UserAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.UserSrv.Update(ctx, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Delete 删除数据
func (a *UserAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Enable 启用数据
func (a *UserAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Disable 禁用数据
func (a *UserAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
