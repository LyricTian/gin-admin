package rbac

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/x/cachex"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type Casbinx struct {
	enforcer               *casbin.Enforcer `wire:"-"`
	Cache                  cachex.Cacher
	RoleRepo               *dao.RoleRepo
	MenuActionResourceRepo *dao.MenuActionResourceRepo
}

func (a *Casbinx) GetEnforcer() *casbin.Enforcer {
	return a.enforcer
}

func (a *Casbinx) Load(ctx context.Context) error {
	roleResult, err := a.RoleRepo.Query(ctx, typed.RoleQueryParam{
		Status: typed.RoleStatusEnabled,
	}, typed.RoleQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, role := range roleResult.Data {
		resourceResult, err := a.MenuActionResourceRepo.Query(ctx, typed.MenuActionResourceQueryParam{
			RoleID: role.ID,
		})
		if err != nil {
			return err
		}
		for _, resource := range resourceResult.Data {
			buf.WriteString(fmt.Sprintf("p, %s, %s, %s \n", role.ID, resource.Path, resource.Method))
		}
	}

	policyFile := filepath.Join(config.C.General.ConfigDir, "rbac_policy.csv")
	_ = os.Rename(policyFile, policyFile+".bak")

	err = ioutil.WriteFile(policyFile, buf.Bytes(), 0666)
	if err != nil {
		logger.Context(ctx).Error("write rbac policy file error", zap.String("file", policyFile), zap.Error(err))
		return err
	}

	// set readonly
	_ = os.Chmod(policyFile, 0444)

	// load casbin
	modelFile := filepath.Join(config.C.General.ConfigDir, "casbin_model.conf")
	e, err := casbin.NewEnforcer(modelFile, policyFile)
	if err != nil {
		return err
	}

	e.EnableEnforce(!config.C.Middleware.Casbin.Disable)
	e.EnableLog(config.C.Middleware.Casbin.Debug)
	a.enforcer = e

	logger.Context(ctx).Info("Load casbin success", zap.String("file", policyFile))

	return nil
}

func (a *Casbinx) AutoLoad(ctx context.Context) {
	updated := time.Now()
	ticker := time.NewTicker(time.Second * time.Duration(config.C.Middleware.Casbin.AutoLoadInterval))
	for range ticker.C {
		// If exists role updated, reload casbin
		roleResult, err := a.RoleRepo.Query(ctx, typed.RoleQueryParam{
			GtUpdatedAt: &updated,
		}, typed.RoleQueryOptions{
			QueryOptions: utilx.QueryOptions{
				OrderFields: []utilx.OrderByParam{
					{Field: "updated_at", Direction: utilx.DESC},
				},
				SelectFields: []string{"updated_at"},
			},
		})
		if err != nil {
			logger.Context(ctx).Error("Failed to query role", zap.Error(err))
			continue
		} else if len(roleResult.Data) > 0 {
			updated = roleResult.Data[0].UpdatedAt
			if err := a.Load(ctx); err != nil {
				logger.Context(ctx).Error("Failed to load casbin", zap.Error(err))
			}
			continue
		}

		// If exists role deleted in cache, reload casbin
		v, ok, err := a.Cache.GetAndDelete(ctx, utilx.CacheNSForDeletedRole, utilx.CacheKeyForDeletedRole)
		if err != nil {
			logger.Context(ctx).Error("Failed to get and delete cache", zap.Error(err))
			continue
		} else if ok && v == "1" {
			updated = time.Now()
			if err := a.Load(ctx); err != nil {
				logger.Context(ctx).Error("Failed to load casbin", zap.Error(err))
			}
			continue
		}

	}
}
