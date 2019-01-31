package auth

import (
	"fmt"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/go-session/session"
)

// NewSessionAuth 创建session授权
func NewSessionAuth(store ...session.ManagerStore) Auther {
	auth := &sessionAuth{
		storeKey: "auth_user_info",
	}
	if len(store) > 0 {
		auth.store = store[0]
	}
	return auth
}

type sessionAuth struct {
	storeKey string
	store    session.ManagerStore
}

func (a *sessionAuth) getFunctionName(name string) string {
	return fmt.Sprintf("auth.session.%s", name)
}

func (a *sessionAuth) Entry(skipper SkipperFunc) gin.HandlerFunc {
	options := config.GetSessionConfig()

	var opts []session.Option
	opts = append(opts, session.SetEnableSetCookie(false))
	opts = append(opts, session.SetEnableSIDInURLQuery(false))
	opts = append(opts, session.SetEnableSIDInHTTPHeader(true))
	opts = append(opts, session.SetSessionNameInHTTPHeader(options.HeaderName))
	opts = append(opts, session.SetSign([]byte(options.Sign)))
	opts = append(opts, session.SetExpired(options.Expired))

	if options.EnableStore && a.store != nil {
		opts = append(opts, session.SetStore(a.store))
	}

	ginConfig := ginsession.DefaultConfig
	ginConfig.Skipper = skipper
	ginConfig.ErrorHandleFunc = func(c *gin.Context, err error) {
		ctx := context.New(c)
		if err == session.ErrInvalidSessionID {
			ctx.ResError(errors.NewBadRequestError("无效的会话"))
			return
		}
		logger.StartSpan(ctx.CContext(), "session中间件", a.getFunctionName("Entry")).Errorf(err.Error())
		ctx.ResError(errors.NewInternalServerError())
	}

	return ginsession.NewWithConfig(ginConfig, opts...)
}

func (a *sessionAuth) SaveUserInfo(c *gin.Context, info UserInfo) error {
	ctx := context.New(c)
	store, err := ginsession.Refresh(c)
	if err != nil {
		logger.StartSpan(ctx.CContext(), "更新会话", a.getFunctionName("SaveUserInfo")).Errorf(err.Error())
		return errors.NewInternalServerError("更新会话发生错误")
	}

	store.Set(a.storeKey, info.String())
	err = store.Save()
	if err != nil {
		logger.StartSpan(ctx.CContext(), "存储会话", a.getFunctionName("SaveUserInfo")).Errorf(err.Error())
		return errors.NewInternalServerError("存储会话发生错误")
	}
	return nil
}

func (a *sessionAuth) GetUserInfo(c *gin.Context) (*UserInfo, error) {
	store := ginsession.FromContext(c)
	userID, ok := store.Get(a.storeKey)
	if !ok || userID == nil {
		return nil, nil
	}
	return parseUserInfo(userID.(string)), nil
}

func (a *sessionAuth) Destroy(c *gin.Context) error {
	ctx := context.New(c)
	err := ginsession.Destroy(c)
	if err != nil {
		logger.StartSpan(ctx.CContext(), "销毁会话", a.getFunctionName("Destroy")).Errorf(err.Error())
		return errors.NewInternalServerError("销毁会话发生错误")
	}
	return nil
}
