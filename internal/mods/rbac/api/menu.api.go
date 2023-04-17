package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/gin-gonic/gin"
)

// Menu management for RBAC
type Menu struct {
	MenuBIZ *biz.Menu
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Query menu list
// @Param name query string false "Display name of menu"
// @Success 200 {object} utils.ResponseResult{data=[]schema.Menu}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := utils.ParseQuery(c, &params); err != nil {
		utils.ResError(c, err)
		return
	}

	result, err := a.MenuBIZ.Query(ctx, params)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResPage(c, result.Data, result.PageResult)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Get menu record by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult{data=schema.Menu}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/menus/{id} [get]
func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, item)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Create menu record
// @Param body body schema.MenuForm true "Request body"
// @Success 200 {object} utils.ResponseResult{data=schema.Menu}
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.MenuForm)
	if err := utils.ParseJSON(c, item); err != nil {
		utils.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		utils.ResError(c, err)
		return
	}

	result, err := a.MenuBIZ.Create(ctx, item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, result)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Update menu record by ID
// @Param id path string true "unique id"
// @Param body body schema.MenuForm true "Request body"
// @Success 200 {object} utils.ResponseResult
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.MenuForm)
	if err := utils.ParseJSON(c, item); err != nil {
		utils.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		utils.ResError(c, err)
		return
	}

	err := a.MenuBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Delete menu record by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/menus/{id} [delete]
func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}
