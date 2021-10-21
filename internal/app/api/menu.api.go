package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/LyricTian/gin-admin/v8/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/internal/app/service"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuAPI), "*"))

type MenuAPI struct {
	MenuSrv *service.MenuSrv
}

func (a *MenuAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.MenuSrv.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(
			schema.NewOrderField("sequence", schema.OrderByDESC),
			schema.NewOrderField("id", schema.OrderByDESC),
		),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *MenuAPI) QueryTree(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.MenuSrv.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, result.Data.ToTree())
}

func (a *MenuAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuSrv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *MenuAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = contextx.FromUserID(ctx)
	result, err := a.MenuSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *MenuAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.MenuSrv.Update(ctx, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *MenuAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *MenuAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *MenuAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
