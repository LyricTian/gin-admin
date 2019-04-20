package ctl

import (
	"fmt"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/bll"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/config"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/ginplus"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/gin-gonic/gin"
)

// NewLogin 创建登录管理控制器
func NewLogin(b *bll.Common, a auth.Auther) *Login {
	return &Login{
		LoginBll: b.Login,
		Auth:     a,
	}
}

// Login 登录管理
// @Name Login
// @Description 登录管理
type Login struct {
	LoginBll *bll.Login
	Auth     auth.Auther
}

func (a *Login) getFuncName(name string) string {
	return fmt.Sprintf("web.ctl.Login.%s", name)
}

// GetCaptchaID 获取验证码ID
// @Summary 获取验证码ID
// @Success 200 schema.LoginCaptcha
// @Router GET /api/v1/login/captchaid
func (a *Login) GetCaptchaID(c *gin.Context) {
	captchaID := captcha.NewLen(config.GetGlobalConfig().Captcha.Length)
	data := schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	ginplus.ResSuccess(c, data)
}

// GetCaptcha 获取图形验证码
// @Summary 获取图形验证码
// @Param id query string true "验证码ID"
// @Param reload query string false "重新加载"
// @Success 200 file "图形验证码"
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/login/captcha
func (a *Login) GetCaptcha(c *gin.Context) {
	captchaID := c.Query("id")
	if captchaID == "" {
		ginplus.ResError(c, errors.NewBadRequestError("无效的请求参数"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginplus.ResError(c, errors.NewBadRequestError("无效的请求参数"))
			return
		}
	}

	w := c.Writer
	captchaConfig := config.GetGlobalConfig().Captcha
	err := captcha.WriteImage(w, captchaID, captchaConfig.Width, captchaConfig.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			ginplus.ResError(c, errors.NewBadRequestError("无效的请求参数"))
			return
		}
		logger.StartSpan(ginplus.NewContext(c), "获取图形验证码", a.getFuncName("GetCaptcha")).Errorf(err.Error())
		ginplus.ResError(c, errors.NewInternalServerError())
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
}

// Login 用户登录
// @Summary 用户登录
// @Param body body schema.LoginParam true
// @Success 200 auth.TokenInfo "{access_token:访问令牌,token_type:令牌类型,expires_in:过期时长(单位秒)}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/login
func (a *Login) Login(c *gin.Context) {
	var item schema.LoginParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginplus.ResError(c, errors.NewBadRequestError("无效的验证码"))
		return
	}

	user, err := a.LoginBll.Verify(ginplus.NewContext(c), item.UserName, item.Password)
	if err != nil {
		switch err {
		case bll.ErrInvalidUserName, bll.ErrInvalidPassword:
			ginplus.ResError(c, errors.NewBadRequestError("用户名或密码错误"))
			return
		case bll.ErrUserDisable:
			ginplus.ResError(c, errors.NewBadRequestError("用户被禁用，请联系管理员"))
			return
		default:
			ginplus.ResError(c, errors.NewInternalServerError())
			return
		}
	}

	userID := user.RecordID
	ginplus.SetUserID(c, userID)
	span := logger.StartSpan(ginplus.NewContext(c), "用户登录", a.getFuncName("Login"))
	tokenInfo, err := a.Auth.GenerateToken(userID)
	if err != nil {
		span.Errorf(err.Error())
		ginplus.ResError(c, errors.NewInternalServerError())
		return
	}

	span.Infof("登入系统")
	ginplus.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
// @Summary 用户登出
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Router POST /api/v1/login/exit
func (a *Login) Logout(c *gin.Context) {
	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := ginplus.GetUserID(c)
	if userID != "" {
		span := logger.StartSpan(ginplus.NewContext(c), "用户登出", a.getFuncName("Logout"))
		err := a.Auth.DestroyToken(ginplus.GetToken(c))
		if err != nil {
			span.Errorf(err.Error())
		}
		span.Infof("登出系统")
	}
	ginplus.ResOK(c)
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 auth.TokenInfo "{access_token:访问令牌,token_type:令牌类型,expires_in:过期时长(单位秒)}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/refresh_token
func (a *Login) RefreshToken(c *gin.Context) {
	tokenInfo, err := a.Auth.GenerateToken(ginplus.GetUserID(c))
	if err != nil {
		logger.StartSpan(ginplus.NewContext(c), "刷新令牌", a.getFuncName("RefreshToken")).Errorf(err.Error())
		ginplus.ResError(c, errors.NewInternalServerError())
		return
	}
	ginplus.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 schema.UserLoginInfo
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/current/user
func (a *Login) GetUserInfo(c *gin.Context) {
	info, err := a.LoginBll.GetUserInfo(ginplus.NewContext(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
// @Summary 查询当前用户菜单树
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 schema.Menu "查询结果：{list:菜单树}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/current/menutree
func (a *Login) QueryUserMenuTree(c *gin.Context) {
	menus, err := a.LoginBll.QueryUserMenuTree(ginplus.NewContext(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, menus)
}

// UpdatePassword 更新个人密码
// @Summary 更新个人密码
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.UpdatePasswordParam true
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/current/password
func (a *Login) UpdatePassword(c *gin.Context) {
	var item schema.UpdatePasswordParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
