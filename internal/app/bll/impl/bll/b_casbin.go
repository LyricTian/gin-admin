package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/casbin/casbin/v2"
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
				logger.Errorf(item.ctx, "The load casbin policy error: %s", err.Error())
			}
		}
	}()
}

// LoadCasbinPolicy 异步加载casbin权限策略
func LoadCasbinPolicy(ctx context.Context, e *casbin.SyncedEnforcer) {
	if len(chCasbinPolicy) > 0 {
		logger.Infof(ctx, "The load casbin policy is already in the wait queue")
		return
	}

	chCasbinPolicy <- &chCasbinPolicyItem{
		ctx: ctx,
		e:   e,
	}
}
