package api

import (
	"gin-admin/src/bll"
	"gin-admin/src/context"
	"gin-admin/src/logger"
	"gin-admin/src/schema"
	"gin-admin/src/util"

	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
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

	nctx := ctx.NewContext()
	userInfo, err := a.LoginBll.Verify(nctx, item.UserName, item.Password)
	if err != nil {
		logger.LoginWithContext(nctx).Errorf("登录发生错误：%s", err.Error())

		status := "error"
		if err == bll.ErrInvalidPassword ||
			err == bll.ErrInvalidUserName ||
			err == bll.ErrUserDisable {
			status = "fail"
		}

		ctx.ResSuccess(gin.H{"status": status})
		return
	}

	// 更新会话
	store, err := ginsession.Refresh(ctx.Context)
	if err != nil {
		logger.LoginWithContext(nctx).Errorf("登录发生错误：%s", err.Error())
		ctx.ResSuccess(gin.H{"status": "error"})
		return
	}

	store.Set(util.SessionKeyUserID, userInfo.RecordID)
	err = store.Save()
	if err != nil {
		logger.LoginWithContext(nctx).Errorf("登录发生错误：%s", err.Error())
		ctx.ResSuccess(gin.H{"status": "error"})
		return
	}
	logger.LoginWithContext(nctx).Infof("登入系统")

	ctx.ResSuccess(gin.H{"status": "success"})
}

// Logout 用户登出
func (a *Login) Logout(ctx *context.Context) {
	nctx := ctx.NewContext()

	userID := ctx.GetUserID()
	if userID != "" {
		store := ginsession.FromContext(ctx.Context)
		err := store.Flush()
		if err != nil {
			logger.LoginWithContext(nctx).Errorf("登出发生错误：%s", err.Error())
			ctx.ResInternalServerError(err)
			return
		}
		logger.LoginWithContext(nctx).Infof("登出系统")
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
	}
	ctx.ResList(menus)
}
