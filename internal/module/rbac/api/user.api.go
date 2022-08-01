package api

import (
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/service"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	UserSvc *service.UserSvc
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Query user list
// @Param status query string false "user status (enabled/disabled)"
// @Param name query string false "name (fuzzy query)"
// @Param username query string false "username (fuzzy query)"
// @Param roleID query string false "role id"
// @Success 200 {object} utilx.ListResult{list=[]typed.User} "query result"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users [get]
func (a *UserAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params typed.UserQueryParam
	if err := utilx.ParseQuery(c, &params); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.UserSvc.Query(ctx, params)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResList(c, result.Data)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Get single user by id
// @Param id path string true "unique id"
// @Success 200 {object} typed.User
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users/{id} [get]
func (a *UserAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserSvc.Get(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Create user
// @Param body body typed.UserCreate true "request body"
// @Success 200 {object} typed.User
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users [post]
func (a *UserAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.UserCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	result, err := a.UserSvc.Create(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Update user by id
// @Param id path int true "unique id"
// @Param body body typed.UserCreate true "request body"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 400 {object} utilx.ErrorResult
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users/{id} [put]
func (a *UserAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.UserCreate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.UserSvc.Update(ctx, c.Param("id"), item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Delete single user by id
// @Param id path string true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users/{id} [delete]
func (a *UserAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSvc.Delete(ctx, c.Param("id"))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Active user by id
// @Param id path int true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users/{id}/active [patch]
func (a *UserAPI) Active(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSvc.UpdateStatus(ctx, c.Param("id"), typed.UserStatusActivated)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Freeze user by id
// @Param id path int true "unique id"
// @Success 200 {object} utilx.OkResult "ok=true"
// @Failure 401 {object} utilx.ErrorResult
// @Failure 500 {object} utilx.ErrorResult
// @Router /api/rbac/v1/users/{id}/freeze [patch]
func (a *UserAPI) Freeze(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSvc.UpdateStatus(ctx, c.Param("id"), typed.UserStatusFreezed)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}
