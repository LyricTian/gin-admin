package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"github.com/gin-gonic/gin"
)

// Menu management for RBAC
type Menu struct {
	MenuBIZ *biz.Menu
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Query menu tree data
// @Param code query string false "Code path of menu (like xxx.xxx.xxx)"
// @Param name query string false "Name of menu"
// @Success 200 {object} util.ResponseResult{data=[]schema.Menu}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.MenuBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Get menu record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult{data=schema.Menu}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/menus/{id} [get]
func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Create menu record
// @Param body body schema.MenuForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=schema.Menu}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.MenuForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.MenuBIZ.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Update menu record by ID
// @Param id path string true "unique id"
// @Param body body schema.MenuForm true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.MenuForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.MenuBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Delete menu record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/menus/{id} [delete]
func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
