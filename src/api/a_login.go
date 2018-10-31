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

	user, err := a.LoginBll.Verify(ctx.NewContext(), item.UserName, item.Password)
	if err != nil {
		logger.LoginWithContext(ctx.NewContext()).Errorf("登录发生错误：%s", err.Error())
		ctx.ResSuccess(gin.H{"status": "error"})
		return
	}

	// 保存会话
	store := ginsession.FromContext(ctx.Context)
	store.Set(util.SessionKeyUserID, user.RecordID)
	err = store.Save()
	if err != nil {
		logger.LoginWithContext(ctx.NewContext()).Errorf("登录发生错误：%s", err.Error())
		ctx.ResSuccess(gin.H{"status": "error"})
		return
	}

	ctx.ResSuccess(gin.H{"status": "success"})
}
