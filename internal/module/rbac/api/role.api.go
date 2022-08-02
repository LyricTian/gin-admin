package api

import (
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/biz"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/gin-gonic/gin"
)

type RoleAPI struct {
	RoleBiz *biz.RoleBiz
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Query role list
// @Param name query string false "role name (fuzzy query)"
// @Param status query string false "role status (enabled/disabled)"
// @Param result query string false " result type (select/paginate, default: paginate)"
// @Success 200 {object} utilx.ListResult{list=[]typed.Role} "query result"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles [get]
func (a *RoleAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params typed.RoleQueryParam
	if err := utilx.ParseQuery(c, &params); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.RoleBiz.Query(ctx, params)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResList(c, result.Data)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Get single role by id
// @Param id path string true "unique id"
// @Success 200 {object} typed.Role
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles/{id} [get]
func (a *RoleAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleBiz.Get(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Create role
// @Param body body typed.RoleCreate true "request body"
// @Success 200 {object} typed.Role
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles [post]
func (a *RoleAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.RoleCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.RoleBiz.Create(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Update role by id
// @Param id path int true "unique id"
// @Param body body typed.RoleCreate true "request body"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles/{id} [put]
func (a *RoleAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.RoleCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.RoleBiz.Update(ctx, c.Param("id"), item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Delete single role by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles/{id} [delete]
func (a *RoleAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBiz.Delete(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Enable role status by id
// @Param id path int true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles/{id}/enable [patch]
func (a *RoleAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBiz.UpdateStatus(ctx, c.Param("id"), typed.RoleStatusEnabled)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Disable role status by id
// @Param id path int true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/roles/{id}/disable [patch]
func (a *RoleAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBiz.UpdateStatus(ctx, c.Param("id"), typed.RoleStatusDisabled)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
