package api

import (
	"github.com/LyricTian/gin-admin/v7/internal/app/bll"
	"github.com/LyricTian/gin-admin/v7/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// MenuSet 注入Menu
var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"))

// Menu 菜单管理
type Menu struct {
	MenuBll *bll.Menu
}

// Query 查询数据
func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.MenuBll.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.MenuBll.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, result.Data.ToTree())
}

// Get 查询指定数据
func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuBll.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// Create 创建数据
func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.MenuBll.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.MenuBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Delete 删除数据
func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuBll.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Enable 启用数据
func (a *Menu) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuBll.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Disable 禁用数据
func (a *Menu) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuBll.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
