package api

import (
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/biz"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/gin-gonic/gin"
)

type LoginAPI struct {
	LoginBiz *biz.LoginBiz
}

// @Tags LoginAPI
// @Summary Get captcha id
// @Success 200 {object} utilx.ResponseResult{data=typed.Captcha}
// @Router /api/rbac/v1/login/captchaid [get]
func (a *LoginAPI) GetCaptchaID(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginBiz.GetCaptchaID(ctx)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, item)
}

// @Tags LoginAPI
// @Summary Write captcha image
// @Param id query string true "CaptchaID"
// @Param reload query string false "Reload captcha image (reload=1)"
// @Produce image/png
// @Success 200 "Captcha image"
// @Failure 404 {object} utilx.ResponseResult
// @Router /api/rbac/v1/login/captcha [get]
func (a *LoginAPI) WriteCaptchaImage(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.LoginBiz.WriteCaptchaImage(ctx, c.Writer, c.Query("id"), c.Query("reload") == "1")
	if err != nil {
		utilx.ResError(c, err)
	}
}

// @Tags LoginAPI
// @Summary Login system by username and password
// @Param body body typed.UserLogin true "Request body"
// @Success 200 {object} utilx.ResponseResult{data=typed.LoginToken}
// @Failure 400 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/login [post]
func (a *LoginAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.UserLogin
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	loginToken, err := a.LoginBiz.Login(ctx, item.TrimSpace())
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, loginToken)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Logout system
// @Success 200 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/current/logout [post]
func (a *LoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.LoginBiz.Logout(ctx, utilx.GetToken(c))
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Refresh current login token
// @Success 200 {object} utilx.ResponseResult{data=typed.LoginToken}
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/current/refreshtoken [post]
func (a *LoginAPI) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginBiz.RefreshToken(ctx)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, tokenInfo)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Get current user
// @Success 200 {object} utilx.ResponseResult{data=typed.User}
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/current/user [get]
func (a *LoginAPI) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginBiz.GetCurrentUser(ctx)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, info)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Update current user login password
// @Param body body typed.LoginPasswordUpdate true "Request body"
// @Success 200 {object} utilx.ResponseResult
// @Failure 400 {object} utilx.ResponseResult
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/current/password [put]
func (a *LoginAPI) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item typed.LoginPasswordUpdate
	if err := utilx.ParseJSON(c, &item); err != nil {
		utilx.ResError(c, err)
		return
	}

	err := a.LoginBiz.UpdatePassword(ctx, item)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResOK(c)
}

// @Tags LoginAPI
// @Security ApiKeyAuth
// @Summary Query current user privilege menu trees
// @Success 200 {object} utilx.ResponseResult{data=[]typed.Menu} "query result"
// @Failure 401 {object} utilx.ResponseResult
// @Failure 500 {object} utilx.ResponseResult
// @Router /api/rbac/v1/current/menus [put]
func (a *LoginAPI) QueryPrivilegeMenus(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := a.LoginBiz.QueryPrivilegeMenus(ctx)
	if err != nil {
		utilx.ResError(c, err)
		return
	}
	utilx.ResSuccess(c, result)
}
