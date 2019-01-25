package ctl

import (
	"fmt"
	"net/http"

	"github.com/LyricTian/gin-admin/src/bll"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/schema"
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
		ctx.ResError(err)
		return
	}

	var result context.HTTPStatus
	recordID, err := a.LoginBll.Verify(ctx.CContext(), item.UserName, item.Password)
	if err != nil {
		result.Status = context.StatusError
		if err == bll.ErrInvalidPassword ||
			err == bll.ErrInvalidUserName ||
			err == bll.ErrUserDisable {
			result.Status = context.StatusFail
		} else {
			logger.StartSpan(ctx.CContext(), "用户登录", "Login").Errorf(err.Error())
		}

		ctx.ResSuccess(result)
		return
	}
	ctx.SetUserID(recordID)

	nctx := ctx.CContext()
	span := logger.StartSpan(nctx, "用户登录", "Login")
	// 更新会话
	store, err := ctx.RefreshSession()
	if err != nil {
		result.Status = context.StatusError
		span.Errorf("更新会话发生错误：%s", err.Error())
		ctx.ResSuccess(result)
		return
	}

	store.Set(context.ContextKeyUserID, recordID)
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
		nctx := ctx.CContext()
		span := logger.StartSpan(nctx, "用户登出", "Logout")
		if err := ctx.DestroySession(); err != nil {
			span.Errorf("登出系统发生错误：%s", err.Error())
			ctx.ResError(err)
			return
		}
		span.Infof("登出系统")
	}

	ctx.ResOK()
}

// GetUserInfo 获取用户登录信息
func (a *Login) GetUserInfo(ctx *context.Context) {
	info, err := a.LoginBll.GetUserInfo(ctx.CContext())
	if err != nil {
		ctx.ResError(err)
		return
	}
	ctx.ResSuccess(info)
}

// QueryCurrentUserMenus 查询当前用户菜单
func (a *Login) QueryCurrentUserMenus(ctx *context.Context) {
	menus, err := a.LoginBll.QueryUserMenuTree(ctx.CContext())
	if err != nil {
		ctx.ResError(err)
		return
	} else if len(menus) == 0 {
		ctx.ResError(fmt.Errorf("用户未授权"), http.StatusUnauthorized, 9998)
		return
	}
	ctx.ResList(menus)
}
