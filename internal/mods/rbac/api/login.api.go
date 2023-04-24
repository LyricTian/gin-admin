package api

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/gin-gonic/gin"
)

type Login struct {
	LoginBIZ *biz.Login
}

// @Tags LoginAPI
// @Summary Get login verify info (captcha id)
// @Success 200 {object} utils.ResponseResult{data=schema.LoginVerify}
// @Router /api/v1/login/verify [get]
func (a *Login) GetVerify(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.GetVerify(ctx)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, data)
}

// @Tags LoginAPI
// @Summary Response captcha image
// @Param id query string true "Captcha ID"
// @Param reload query number false "Reload captcha image (reload=1)"
// @Produce image/png
// @Success 200 "Captcha image"
// @Failure 404 {object} utils.ResponseResult
// @Router /api/v1/login/captcha [get]
func (a *Login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.LoginBIZ.ResCaptcha(ctx, c.Writer, c.Query("id"), c.Query("reload") == "1")
	if err != nil {
		utils.ResError(c, err)
	}
}

// @Tags LoginAPI
// @Summary Login system with username and password
// @Param body body schema.LoginForm true "Request body"
// @Success 200 {object} utils.ResponseResult{data=schema.LoginToken}
// @Failure 400 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/login [post]
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.LoginForm)
	if err := utils.ParseJSON(c, item); err != nil {
		utils.ResError(c, err)
		return
	}

	data, err := a.LoginBIZ.Login(ctx, item.Trim())
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, data)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Logout system
// @Success 200 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/current/logout [post]
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.LoginBIZ.Logout(ctx)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Refresh current access token
// @Success 200 {object} utils.ResponseResult{data=schema.LoginToken}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/current/refresh-token [post]
func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.RefreshToken(ctx)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, data)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Get current user info
// @Success 200 {object} utils.ResponseResult{data=schema.User}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/current/user [get]
func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.GetUserInfo(ctx)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, data)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Change current user password
// @Param body body schema.UpdateLoginPassword true "Request body"
// @Success 200 {object} utils.ResponseResult
// @Failure 400 {object} utils.ResponseResult
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/current/password [put]
func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UpdateLoginPassword)
	if err := utils.ParseJSON(c, item); err != nil {
		utils.ResError(c, err)
		return
	}

	err := a.LoginBIZ.UpdatePassword(ctx, item)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResOK(c)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Query current user menus based on the current user role
// @Success 200 {object} utils.ResponseResult{data=[]schema.Menu}
// @Failure 401 {object} utils.ResponseResult
// @Failure 500 {object} utils.ResponseResult
// @Router /api/v1/current/menus [get]
func (a *Login) QueryMenus(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.QueryMenus(ctx)
	if err != nil {
		utils.ResError(c, err)
		return
	}
	utils.ResSuccess(c, data)
}
