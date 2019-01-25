package ctl

import (
	"fmt"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Login 登录管理
// @Name Login
// @Description 登录管理
type Login struct {
	LoginBll *bll.Login `inject:""`
}

func (a *Login) getFunctionName(name string) string {
	return fmt.Sprintf("web.ctl.Login.%s", name)
}

// GetCaptchaID 获取验证码ID
// @Summary 获取验证码ID
// @Success 200 schema.LoginCaptcha
// @Router GET /api/v1/login/captchaid
func (a *Login) GetCaptchaID(ctx *context.Context) {
	captchaID := captcha.NewLen(config.GetCaptchaConfig().Length)
	data := schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	ctx.ResSuccess(data)
}

// GetCaptcha 获取图形验证码
// @Summary 获取图形验证码
// @Param id query string true "验证码ID"
// @Param reload query int false "是否重新加载（1是）"
// @Success 200 file "图形验证码"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/login/captcha
func (a *Login) GetCaptcha(ctx *context.Context) {
	captchaID := ctx.Query("id")
	if captchaID == "" {
		ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
		return
	}

	if ctx.Query("reload") == "1" {
		if !captcha.Reload(captchaID) {
			ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
			return
		}
	}

	w := ctx.ResponseWriter()
	captchaConfig := config.GetCaptchaConfig()
	err := captcha.WriteImage(w, captchaID, captchaConfig.Width, captchaConfig.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			ctx.ResError(errors.NewBadRequestError("无效的请求参数"))
			return
		}
		logger.StartSpan(ctx.CContext(), "获取图形验证码", a.getFunctionName("GetCaptcha")).Errorf(err.Error())
		ctx.ResError(errors.NewInternalServerError())
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
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/login
func (a *Login) Login(ctx *context.Context) {
	var item schema.LoginParam
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	if !captcha.VerifyString(item.VerifyID, item.VerifyCode) {
		ctx.ResError(errors.NewBadRequestError("无效的验证码"))
		return
	}

	userID, err := a.LoginBll.Verify(ctx.CContext(), item.UserName, item.Password)
	if err != nil {
		switch err {
		case bll.ErrInvalidUserName, bll.ErrInvalidPassword:
			ctx.ResError(errors.NewBadRequestError("用户名或密码错误"))
			return
		case bll.ErrUserDisable:
			ctx.ResError(errors.NewBadRequestError("用户被禁用，请联系管理员"))
			return
		default:
			ctx.ResError(errors.NewInternalServerError())
			return
		}
	}
	ctx.SetUserID(userID)

	err = ctx.SaveUserIDToSession(userID)
	if err != nil {
		ctx.ResError(err)
		return
	}

	logger.StartSpan(ctx.CContext(), "用户登录", a.getFunctionName("Login")).Infof("登入系统")
	ctx.ResOK()
}

// Logout 用户登出
// @Summary 用户登出
// @Success 200 option.Interface "{status:OK}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/login/exit
func (a *Login) Logout(ctx *context.Context) {
	// 检查用户是否处于登录状态，如果是，则执行销毁
	userID := ctx.GetUserID()
	if userID != "" {
		if err := ctx.DestroySession(); err != nil {
			ctx.ResError(err)
			return
		}
		logger.StartSpan(ctx.CContext(), "用户登出", a.getFunctionName("Logout")).Infof("登出系统")
	}
	ctx.ResOK()
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Param Access-Token header string false "访问令牌"
// @Success 200 schema.UserLoginInfo
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/current/user
func (a *Login) GetUserInfo(ctx *context.Context) {
	info, err := a.LoginBll.GetUserInfo(ctx.CContext())
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(info)
}

// QueryUserMenuTree 查询当前用户菜单树
// @Summary 查询当前用户菜单树
// @Param Access-Token header string false "访问令牌"
// @Success 200 option.Interface "查询结果：{list:菜单树}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/current/menutree
func (a *Login) QueryUserMenuTree(ctx *context.Context) {
	menus, err := a.LoginBll.QueryUserMenuTree(ctx.CContext())
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResList(menus)
}

// UpdatePassword 更新个人密码
// @Summary 更新个人密码
// @Param body body schema.UpdatePasswordParam true
// @Success 200 option.Interface "{status:OK}"
// @Failure 400 option.Interface "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 option.Interface "{error:{code:0,message:未授权}}"
// @Failure 500 option.Interface "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/current/password
func (a *Login) UpdatePassword(ctx *context.Context) {
	var item schema.UpdatePasswordParam
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResError(err)
		return
	}

	err := a.LoginBll.UpdatePassword(ctx.CContext(), item)
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResOK()
}
