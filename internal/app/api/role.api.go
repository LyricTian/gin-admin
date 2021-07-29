package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/LyricTian/gin-admin/v8/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/internal/app/service"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleAPI), "*"))

// RoleAPI 角色管理
type RoleAPI struct {
	RoleSrv *service.RoleSrv
}

// Query 查询数据
func (a *RoleAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.RoleSrv.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// QuerySelect 查询选择数据
func (a *RoleAPI) QuerySelect(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.RoleSrv.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, result.Data)
}

// Get 查询指定数据
func (a *RoleAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleSrv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// Create 创建数据
func (a *RoleAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = contextx.FromUserID(ctx)
	result, err := a.RoleSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
func (a *RoleAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.RoleSrv.Update(ctx, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Delete 删除数据
func (a *RoleAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Enable 启用数据
func (a *RoleAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Disable 禁用数据
func (a *RoleAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
