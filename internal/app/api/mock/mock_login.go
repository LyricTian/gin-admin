package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// LoginSet 注入Login
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

// Login 登录管理
type Login struct {
}

// GetCaptcha 获取验证码信息
// @Tags 登录管理
// @Summary 获取验证码信息
// @Success 200 {object} schema.LoginCaptcha
// @Router /api/v1/pub/login/captchaid [get]
func (a *Login) GetCaptcha(c *gin.Context) {
}

// ResCaptcha 响应图形验证码
// @Tags 登录管理
// @Summary 响应图形验证码
// @Param id query string true "验证码ID"
// @Param reload query string false "重新加载"
// @Produce image/png
// @Success 200 "图形验证码"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/login/captcha [get]
func (a *Login) ResCaptcha(c *gin.Context) {
}

// Login 用户登录
// @Tags 登录管理
// @Summary 用户登录
// @Param body body schema.LoginParam true "请求参数"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/login [post]
func (a *Login) Login(c *gin.Context) {
}

// Logout 用户登出
// @Tags 登录管理
// @Summary 用户登出
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Router /api/v1/pub/login/exit [post]
func (a *Login) Logout(c *gin.Context) {
}

// RefreshToken 刷新令牌
// @Tags 登录管理
// @Summary 刷新令牌
// @Security ApiKeyAuth
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/refresh-token [post]
func (a *Login) RefreshToken(c *gin.Context) {
}

// GetUserInfo 获取当前用户信息
// @Tags 登录管理
// @Summary 获取当前用户信息
// @Security ApiKeyAuth
// @Success 200 {object} schema.UserLoginInfo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/user [get]
func (a *Login) GetUserInfo(c *gin.Context) {
}

// QueryUserMenuTree 查询当前用户菜单树
// @Tags 登录管理
// @Summary 查询当前用户菜单树
// @Security ApiKeyAuth
// @Success 200 {object} schema.Menu "查询结果：{list:菜单树}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/menutree [get]
func (a *Login) QueryUserMenuTree(c *gin.Context) {
}

// UpdatePassword 更新个人密码
// @Tags 登录管理
// @Summary 更新个人密码
// @Security ApiKeyAuth
// @Param body body schema.UpdatePasswordParam true "请求参数"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/password [put]
func (a *Login) UpdatePassword(c *gin.Context) {
}
