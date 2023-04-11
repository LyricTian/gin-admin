package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/library/utils"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/gin-gonic/gin"
)

type Resource struct {
	ResourceBIZ *biz.Resource
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Query paginated resource list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param code query string false "resource code (fuzzy query)"
// @Param status query string false "resource status (enabled, disabled)"
// @Success 200 {object} utils.ResponseResult{data=[]schema.Resource}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/resources [get]
func (a *Resource) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.ResourceQueryParam
	if err := utils.ParseQuery(c, &params); err != nil {
		utils.ResError(c, err)
		return
	}

	result, err := a.ResourceBIZ.Query(ctx, params)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResPage(c, result.Data, result.PageResult)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Get resource details by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult{data=schema.Resource}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/resources/{id} [get]
func (a *Resource) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.ResourceBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, item)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Create resource record
// @Param body body schema.ResourceSave true "Request body"
// @Success 200 {object} utils.ResponseResult{data=schema.Resource}
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/resources [post]
func (a *Resource) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.ResourceSave
	if err := utils.ParseJSON(c, &item); err != nil {
		utils.ResError(c, err)
		return
	}

	result, err := a.ResourceBIZ.Create(ctx, item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, result)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Update resource record by ID
// @Param id path string true "unique id"
// @Param body body schema.ResourceSave true "Request body"
// @Success 200 {object} utils.ResponseResult
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/resources/{id} [put]
func (a *Resource) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.ResourceSave
	if err := utils.ParseJSON(c, &item); err != nil {
		utils.ResError(c, err)
		return
	}

	err := a.ResourceBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}

// @Tags ResourceAPI
// @Security ApiKeyAuth
// @Summary Delete resource record by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/resources/{id} [delete]
func (a *Resource) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.ResourceBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}
