package api

import (
	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/v7/internal/app/bll"
	"github.com/LyricTian/gin-admin/v7/internal/app/config"
	"github.com/LyricTian/gin-admin/v7/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/errors"
	"github.com/LyricTian/gin-admin/v7/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// LoginSet 注入Login
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

// Login 登录管理
type Login struct {
	LoginBll *bll.Login
}

// GetCaptcha 获取验证码信息
func (a *Login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginBll.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// ResCaptcha 响应图形验证码
func (a *Login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		ginx.ResError(c, errors.New400Response("请提供验证码ID"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginx.ResError(c, errors.New400Response("未找到验证码ID"))
			return
		}
	}

	cfg := config.C.Captcha
	err := a.LoginBll.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginx.ResError(c, err)
	}
}

// Login 用户登录
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginx.ResError(c, errors.New400Response("无效的验证码"))
		return
	}

	user, err := a.LoginBll.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	userID := user.ID
	// 将用户ID放入上下文
	ginx.SetUserID(c, userID)

	tokenInfo, err := a.LoginBll.GenerateToken(ctx, userID)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("登入系统")
	ginx.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := ginx.GetUserID(c)
	if userID != "" {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.LoginBll.DestroyToken(ctx, ginx.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("登出系统")
	}
	ginx.ResOK(c)
}

// RefreshToken 刷新令牌
func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginBll.GetLoginInfo(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
func (a *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginBll.QueryUserMenuTree(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, menus)
}

// UpdatePassword 更新个人密码
func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ctx, ginx.GetUserID(c), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
