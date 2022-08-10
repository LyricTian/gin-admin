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
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param name query string false "role name (fuzzy query)"
// @Param status query string false "role status (enabled/disabled)"
// @Param result query string false " result type (select/paginate, default: paginate)"
// @Success 200 {object} utilx.ResponseResult{data=[]typed.Role} "query result"
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
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
	utilx.ResPage(c, result.Data, result.PageResult)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Get single role by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.ResponseResult{data=typed.Role}
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
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
// @Success 200 {object} utilx.ResponseResult{data=typed.Role}
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
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
// @Success 200 {object} utilx.ResponseResult
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
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
// @Summary Delete role by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
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
// @Summary Update role status by id
// @Param id path int true "unique id"
// @Param body body typed.RoleUpdateStatus true "request body"
// @Success 200 {object} utilx.ResponseResult
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/roles/{id}/status [put]
func (a *RoleAPI) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.RoleUpdateStatus
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.RoleBiz.UpdateStatus(ctx, c.Param("id"), item.Status)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
