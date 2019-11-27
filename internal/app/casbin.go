package app

import (
	"context"
	"fmt"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"go.uber.org/dig"
)

// NewCasbinEnforcer 创建casbin校验器
func NewCasbinEnforcer() *casbin.SyncedEnforcer {
	cfg := config.Global().Casbin
	if !cfg.Enable {
		return nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	handleError(err)

	e.EnableAutoSave(false)
	e.EnableAutoBuildRoleLinks(true)

	if cfg.Debug {
		e.EnableLog(true)
	}
	return e
}

// InitCasbinEnforcer 初始化casbin校验器
func InitCasbinEnforcer(container *dig.Container) error {
	cfg := config.Global().Casbin
	if !cfg.Enable {
		return nil
	}

	return container.Invoke(func(e *casbin.SyncedEnforcer, bRole bll.IRole, bUser bll.IUser) error {
		adapter := NewCasbinAdapter(bRole, bUser)

		if cfg.AutoLoad {
			_ = e.InitWithModelAndAdapter(e.GetModel(), adapter)
			e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		} else {
			err := adapter.LoadPolicy(e.GetModel())
			if err != nil {
				return err
			}
		}

		err := e.BuildRoleLinks()
		if err != nil {
			return err
		}

		return nil
	})
}

// ReleaseCasbinEnforcer 释放casbin资源
func ReleaseCasbinEnforcer(container *dig.Container) {
	cfg := config.Global().Casbin
	if !cfg.Enable || !cfg.AutoLoad {
		return
	}

	_ = container.Invoke(func(e *casbin.SyncedEnforcer) {
		e.StopAutoLoadPolicy()
	})
}

// NewCasbinAdapter 创建casbin适配器
func NewCasbinAdapter(bRole bll.IRole, bUser bll.IUser) *CasbinAdapter {
	return &CasbinAdapter{
		RoleBll: bRole,
		UserBll: bUser,
	}
}

// CasbinAdapter casbin适配器
type CasbinAdapter struct {
	RoleBll bll.IRole
	UserBll bll.IUser
}

// LoadPolicy loads all policy rules from the storage.
func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		logger.Errorf(ctx, "Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.loadUserPolicy(ctx, model)
	if err != nil {
		logger.Errorf(ctx, "Load casbin user policy error: %s", err.Error())
		return err
	}
	return nil
}

func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, model model.Model) error {
	// 加载角色策略
	roleResult, err := a.RoleBll.Query(ctx, schema.RoleQueryParam{}, schema.RoleQueryOptions{
		IncludeMenus: true,
	})
	if err != nil {
		return err
	}

	for _, item := range roleResult.Data {
		resources, err := a.RoleBll.GetMenuResources(ctx, item)
		if err != nil {
			return err
		}

		for _, ritem := range resources {
			if ritem.Path == "" || ritem.Method == "" {
				continue
			}

			line := fmt.Sprintf("p,%s,%s,%s", item.RecordID, ritem.Path, ritem.Method)
			persist.LoadPolicyLine(line, model)
		}
	}

	return nil
}

func (a *CasbinAdapter) loadUserPolicy(ctx context.Context, model model.Model) error {
	result, err := a.UserBll.Query(ctx, schema.UserQueryParam{
		Status: 1,
	}, schema.UserQueryOptions{IncludeRoles: true})
	if err != nil {
		return err
	}

	for _, item := range result.Data {
		for _, roleID := range item.Roles.ToRoleIDs() {
			line := fmt.Sprintf("g,%s,%s", item.RecordID, roleID)
			persist.LoadPolicyLine(line, model)
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
