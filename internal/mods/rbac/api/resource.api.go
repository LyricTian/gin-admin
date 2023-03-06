package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/gin-gonic/gin"
)

type Resource struct {
	ResourceBiz *biz.Resource
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Query paginated resource list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param code query string false "resource code (fuzzy query)"
// @Param status query string false "resource status (enabled, disabled)"
// @Success 200 {object} utilx.ResponseResult{data=[]schema.Resource}
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/v1/resources [get]
func (a *Resource) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.ResourceQueryParam
	if err := utilx.ParseQuery(c, &params); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.ResourceBiz.Query(ctx, params)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResPage(c, result.Data, result.PageResult)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Get resource details by ID
// @Param id path string true "unique id"
// @Success 200 {object} utilx.ResponseResult{data=schema.Resource}
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/v1/resources/{id} [get]
func (a *Resource) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.ResourceBiz.Get(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Create resource record
// @Param body body schema.ResourceCreate true "Request body"
// @Success 200 {object} utilx.ResponseResult{data=schema.Resource}
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/v1/resources [post]
func (a *Resource) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.ResourceCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.ResourceBiz.Create(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Update resource record by ID
// @Param id path string true "unique id"
// @Param body body schema.ResourceCreate true "Request body"
// @Success 200 {object} utilx.ResponseResult
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/v1/resources/{id} [put]
func (a *Resource) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.ResourceCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.ResourceBiz.Update(ctx, c.Param("id"), item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Delete resource record by ID
// @Param id path string true "unique id"
// @Success 200 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/v1/resources/{id} [delete]
func (a *Resource) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.ResourceBiz.Delete(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
