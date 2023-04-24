package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/gin-gonic/gin"
)

// User management for RBAC
type User struct {
	UserBIZ *biz.User
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Query user list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param username query string false "Username for login"
// @Param name query string false "Name of user"
// @Param status query string false "Status of user (activated, freezed)"
// @Success 200 {object} utils.ResponseResult{data=[]schema.User}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := utils.ParseQuery(c, &params); err != nil {
		utils.ResError(c, err)
		return
	}

	result, err := a.UserBIZ.Query(ctx, params)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResPage(c, result.Data, result.PageResult)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Get user record by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult{data=schema.User}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, item)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Create user record
// @Param body body schema.UserForm true "Request body"
// @Success 200 {object} utils.ResponseResult{data=schema.User}
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UserForm)
	if err := utils.ParseJSON(c, item); err != nil {
		utils.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		utils.ResError(c, err)
		return
	}

	result, err := a.UserBIZ.Create(ctx, item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, result)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Update user record by ID
// @Param id path string true "unique id"
// @Param body body schema.UserForm true "Request body"
// @Success 200 {object} utils.ResponseResult
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UserForm)
	if err := utils.ParseJSON(c, item); err != nil {
		utils.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		utils.ResError(c, err)
		return
	}

	err := a.UserBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Delete user record by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Reset user password by ID
// @Param id path string true "unique id"
// @Success 200 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/users/{id}/reset-pwd [patch]
func (a *User) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBIZ.ResetPassword(ctx, c.Param("id"))
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}
