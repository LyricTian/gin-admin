package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v9/internal/module/ginx"
	"github.com/LyricTian/gin-admin/v9/internal/schema"
	"github.com/LyricTian/gin-admin/v9/internal/service"
)

var DemoSet = wire.NewSet(wire.Struct(new(DemoAPI), "*"))

type DemoAPI struct {
	DemoSrv *service.DemoSrv
}

// @Tags DemoAPI
// @Summary Query demo list
// @Security ApiKeyAuth
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Success 200 {object} schema.ListResult{list=[]schema.Demo} "Query result (schema.Demo object)"
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/demos [get]
func (a *DemoAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.DemoQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.DemoSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// @Tags DemoAPI
// @Summary Get single demo by id
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Success 200 {object} schema.Demo
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/demos/{id} [get]
func (a *DemoAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.DemoSrv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// @Tags DemoAPI
// @Summary Create demo
// @Security ApiKeyAuth
// @Param body body schema.Demo true "Request body"
// @Success 200 {object} schema.Demo
// @Failure 400 {object} schema.ErrorResult
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/demos [post]
func (a *DemoAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Demo
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.DemoSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// @Tags DemoAPI
// @Summary Update demo by id
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Param body body schema.Demo true "Request body"
// @Success 200 {object} schema.OkResult "ok=true"
// @Failure 400 {object} schema.ErrorResult
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/demos/{id} [put]
func (a *DemoAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Demo
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.DemoSrv.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// @Tags DemoAPI
// @Summary Delete single demo by id
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Success 200 {object} schema.OkResult "ok=true"
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/demos/{id} [delete]
func (a *DemoAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.DemoSrv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
