package service

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
)

// ServiceSet bll注入
var ServiceSet = wire.NewSet(
	MenuSet,
	RoleSet,
	UserSet,
	LoginSet,
)

var chCasbinPolicy chan *chCasbinPolicyItem

type chCasbinPolicyItem struct {
	ctx context.Context
	e   *casbin.SyncedEnforcer
}

func init() {
	chCasbinPolicy = make(chan *chCasbinPolicyItem, 1)
	go func() {
		for item := range chCasbinPolicy {
			err := item.e.LoadPolicy()
			if err != nil {
				logger.WithContext(item.ctx).Errorf("The load casbin policy error: %s", err.Error())
			}
		}
	}()
}

// LoadCasbinPolicy 异步加载casbin权限策略
func LoadCasbinPolicy(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !config.C.Casbin.Enable || config.C.Casbin.AutoLoad {
		return
	}

	if len(chCasbinPolicy) > 0 {
		logger.WithContext(ctx).Infof("The load casbin policy is already in the wait queue")
		return
	}

	chCasbinPolicy <- &chCasbinPolicyItem{
		ctx: ctx,
		e:   e,
	}
}
