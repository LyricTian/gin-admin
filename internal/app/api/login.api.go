package api

import (
	"fmt"

	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/LyricTian/gin-admin/v8/internal/app/ginx"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/internal/app/service"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
)

var LoginSet = wire.NewSet(wire.Struct(new(LoginAPI), "*"))

// LoginAPI 登录管理
type LoginAPI struct {
	LoginSrv *service.LoginSrv
}

// GetCaptcha 获取验证码信息
func (a *LoginAPI) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginSrv.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// ResCaptcha 响应图形验证码
func (a *LoginAPI) ResCaptcha(c *gin.Context) {
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
	err := a.LoginSrv.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginx.ResError(c, err)
	}
}

// Login 用户登录
func (a *LoginAPI) Login(c *gin.Context) {
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

	user, err := a.LoginSrv.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, a.formatTokenUserID(user.ID, user.UserName))
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, user.ID)
	ctx = logger.NewUserNameContext(ctx, user.UserName)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("登入系统")

	ginx.ResSuccess(c, tokenInfo)
}

func (a *LoginAPI) formatTokenUserID(userID uint64, userName string) string {
	return fmt.Sprintf("%d-%s", userID, userName)
}

// Logout 用户登出
func (a *LoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := contextx.FromUserID(ctx)
	if userID != 0 {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.LoginSrv.DestroyToken(ctx, ginx.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("登出系统")
	}
	ginx.ResOK(c)
}

// RefreshToken 刷新令牌
func (a *LoginAPI) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, a.formatTokenUserID(contextx.FromUserID(ctx), contextx.FromUserName(ctx)))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
func (a *LoginAPI) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginSrv.GetLoginInfo(ctx, contextx.FromUserID(ctx))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
func (a *LoginAPI) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginSrv.QueryUserMenuTree(ctx, contextx.FromUserID(ctx))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, menus)
}

// UpdatePassword 更新个人密码
func (a *LoginAPI) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.LoginSrv.UpdatePassword(ctx, contextx.FromUserID(ctx), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
