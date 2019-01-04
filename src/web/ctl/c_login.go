package ctl

import (
	"fmt"
	"net/http"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
)

// Login 登录管理
type Login struct {
	LoginBll *bll.Login `inject:""`
}

// Login 用户登录
func (a *Login) Login(ctx *context.Context) {
	var item schema.LoginParam
	if err := ctx.ParseJSON(&item); err != nil {
		ctx.ResBadRequest(err)
		return
	}

	var result context.HTTPStatus
	userInfo, err := a.LoginBll.Verify(ctx.NewContext(), item.UserName, item.Password)
	if err != nil {
		result.Status = context.StatusError
		if err == bll.ErrInvalidPassword ||
			err == bll.ErrInvalidUserName ||
			err == bll.ErrUserDisable {
			result.Status = context.StatusFail
		} else {
			logger.Start(ctx.NewContext()).Errorf("用户登录发生错误：%s", err.Error())
		}

		ctx.ResSuccess(result)
		return
	}
	ctx.SetUserID(userInfo.RecordID)

	nctx := ctx.NewContext()
	span := logger.StartSpan(nctx, "用户登录", "Login")
	// 更新会话
	store, err := ctx.RefreshSession()
	if err != nil {
		result.Status = context.StatusError
		span.Errorf("更新会话发生错误：%s", err.Error())
		ctx.ResSuccess(result)
		return
	}

	store.Set(util.SessionKeyUserID, userInfo.RecordID)
	err = store.Save()
	if err != nil {
		result.Status = context.StatusError
		span.Errorf("存储会话发生错误：%s", err.Error())
		ctx.ResSuccess(result)
		return
	}
	span.Infof("登入系统")

	ctx.ResOK()
}

// Logout 用户登出
func (a *Login) Logout(ctx *context.Context) {
	userID := ctx.GetUserID()
	if userID != "" {
		nctx := ctx.NewContext()
		span := logger.StartSpan(nctx, "用户登出", "Logout")
		if err := ctx.DestroySession(); err != nil {
			span.Errorf("登出系统发生错误：%s", err.Error())
			ctx.ResInternalServerError(err)
			return
		}
		span.Infof("登出系统")
	}

	ctx.ResOK()
}

// GetCurrentUserInfo 获取当前用户信息
func (a *Login) GetCurrentUserInfo(ctx *context.Context) {
	userID := ctx.GetUserID()

	info, err := a.LoginBll.GetCurrentUserInfo(ctx.NewContext(), userID)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	}
	ctx.ResSuccess(info)
}

// QueryCurrentUserMenus 查询当前用户菜单
func (a *Login) QueryCurrentUserMenus(ctx *context.Context) {
	userID := ctx.GetUserID()

	menus, err := a.LoginBll.QueryCurrentUserMenus(ctx.NewContext(), userID)
	if err != nil {
		ctx.ResInternalServerError(err)
		return
	} else if len(menus) == 0 {
		ctx.ResError(fmt.Errorf("用户未授权"), http.StatusUnauthorized, 9998)
		return
	}
	ctx.ResList(menus)
}
