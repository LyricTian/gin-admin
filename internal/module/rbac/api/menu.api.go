package api

import (
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/service"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"

	"github.com/gin-gonic/gin"
)

type MenuAPI struct {
	MenuSvc *service.MenuSvc
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Query menu tree
// @Param parentID query int false "parent id"
// @Param name query string false "menu name (fuzzy query)"
// @Param status query string false "menu status (enabled/disabled)"
// @Success 200 {object} utilx.ListResult{list=[]typed.Menu} "query result"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus [get]
func (a *MenuAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params typed.MenuQueryParam
	if err := utilx.ParseQuery(c, &params); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.MenuSvc.Query(ctx, params)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResList(c, result.Data)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Get single menu by id
// @Param id path string true "unique id"
// @Success 200 {object} typed.Menu
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus/{id} [get]
func (a *MenuAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuSvc.Get(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Create menu
// @Param body body typed.MenuCreate true "request body"
// @Success 200 {object} typed.Menu
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus [post]
func (a *MenuAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.MenuCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.MenuSvc.Create(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Update menu by id
// @Param id path int true "unique id"
// @Param body body typed.MenuCreate true "request body"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus/{id} [put]
func (a *MenuAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.MenuCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.MenuSvc.Update(ctx, c.Param("id"), item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Delete single menu by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus/{id} [delete]
func (a *MenuAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSvc.Delete(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Enable menu status by id
// @Param id path int true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus/{id}/enable [patch]
func (a *MenuAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSvc.UpdateStatus(ctx, c.Param("id"), typed.MenuStatusEnabled)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Disable menu status by id
// @Param id path int true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/menus/{id}/disable [patch]
func (a *MenuAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSvc.UpdateStatus(ctx, c.Param("id"), typed.MenuStatusDisabled)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
