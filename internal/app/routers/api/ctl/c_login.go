package ctl

import (
	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/ginplus"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/gin-gonic/gin"
)

// NewLogin 创建登录管理控制器
func NewLogin(bLogin bll.ILogin) *Login {
	return &Login{
		LoginBll: bLogin,
	}
}

// Login 登录管理
type Login struct {
	LoginBll bll.ILogin
}

// GetCaptcha 获取验证码信息
// @Tags 登录管理
// @Summary 获取验证码信息
// @Success 200 {object} schema.LoginCaptcha
// @Router /api/v1/pub/login/captchaid [get]
func (a *Login) GetCaptcha(c *gin.Context) {
	item, err := a.LoginBll.GetCaptcha(ginplus.NewContext(c), config.Global().Captcha.Length)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// ResCaptcha 响应图形验证码
// @Tags 登录管理
// @Summary 响应图形验证码
// @Param id query string true "验证码ID"
// @Param reload query string false "重新加载"
// @Produce image/png
// @Success 200 "图形验证码"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/login/captcha [get]
func (a *Login) ResCaptcha(c *gin.Context) {
	captchaID := c.Query("id")
	if captchaID == "" {
		ginplus.ResError(c, errors.New400Response("请提供验证码ID"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginplus.ResError(c, errors.New400Response("未找到验证码ID"))
			return
		}
	}

	cfg := config.Global().Captcha
	err := a.LoginBll.ResCaptcha(ginplus.NewContext(c), c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginplus.ResError(c, err)
	}
}

// Login 用户登录
// @Tags 登录管理
// @Summary 用户登录
// @Param body body schema.LoginParam true "请求参数"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/login [post]
func (a *Login) Login(c *gin.Context) {
	var item schema.LoginParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginplus.ResError(c, errors.New400Response("无效的验证码"))
		return
	}

	user, err := a.LoginBll.Verify(ginplus.NewContext(c), item.UserName, item.Password)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	userID := user.RecordID
	// 将用户ID放入上下文
	ginplus.SetUserID(c, userID)

	tokenInfo, err := a.LoginBll.GenerateToken(ginplus.NewContext(c), userID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	logger.StartSpan(ginplus.NewContext(c), logger.SetSpanTitle("用户登录"), logger.SetSpanFuncName("Login")).Infof("登入系统")
	ginplus.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
// @Tags 登录管理
// @Summary 用户登出
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Router /api/v1/pub/login/exit [post]
func (a *Login) Logout(c *gin.Context) {
	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := ginplus.GetUserID(c)
	if userID != "" {
		ctx := ginplus.NewContext(c)
		err := a.LoginBll.DestroyToken(ctx, ginplus.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(ginplus.NewContext(c), logger.SetSpanTitle("用户登出"), logger.SetSpanFuncName("Logout")).Infof("登出系统")
	}
	ginplus.ResOK(c)
}

// RefreshToken 刷新令牌
// @Tags 登录管理
// @Summary 刷新令牌
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/refresh-token [post]
func (a *Login) RefreshToken(c *gin.Context) {
	tokenInfo, err := a.LoginBll.GenerateToken(ginplus.NewContext(c), ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
// @Tags 登录管理
// @Summary 获取当前用户信息
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} schema.UserLoginInfo
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/user [get]
func (a *Login) GetUserInfo(c *gin.Context) {
	info, err := a.LoginBll.GetLoginInfo(ginplus.NewContext(c), ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
// @Tags 登录管理
// @Summary 查询当前用户菜单树
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} schema.Menu "查询结果：{list:菜单树}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/menutree [get]
func (a *Login) QueryUserMenuTree(c *gin.Context) {
	menus, err := a.LoginBll.QueryUserMenuTree(ginplus.NewContext(c), ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, menus)
}

// UpdatePassword 更新个人密码
// @Tags 登录管理
// @Summary 更新个人密码
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.UpdatePasswordParam true "请求参数"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/password [get]
func (a *Login) UpdatePassword(c *gin.Context) {
	var item schema.UpdatePasswordParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ginplus.NewContext(c), ginplus.GetUserID(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
