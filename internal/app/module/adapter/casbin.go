package adapter

import (
	"context"
	"fmt"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
)

var _ persist.Adapter = (*CasbinAdapter)(nil)

var CasbinAdapterSet = wire.NewSet(wire.Struct(new(CasbinAdapter), "*"), wire.Bind(new(persist.Adapter), new(*CasbinAdapter)))

type CasbinAdapter struct {
	RoleRepo         *dao.RoleRepo
	RoleMenuRepo     *dao.RoleMenuRepo
	MenuResourceRepo *dao.MenuActionResourceRepo
	UserRepo         *dao.UserRepo
	UserRoleRepo     *dao.UserRoleRepo
}

// Loads all policy rules from the storage.
func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.loadUserPolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}

	return nil
}

// Load role policy (p,role_id,path,method)
func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	roleResult, err := a.RoleRepo.Query(ctx, schema.RoleQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(roleResult.Data) == 0 {
		return nil
	}

	roleMenuResult, err := a.RoleMenuRepo.Query(ctx, schema.RoleMenuQueryParam{})
	if err != nil {
		return err
	}
	mRoleMenus := roleMenuResult.Data.ToRoleIDMap()

	menuResourceResult, err := a.MenuResourceRepo.Query(ctx, schema.MenuActionResourceQueryParam{})
	if err != nil {
		return err
	}
	mMenuResources := menuResourceResult.Data.ToActionIDMap()

	for _, item := range roleResult.Data {
		mcache := make(map[string]struct{})
		if rms, ok := mRoleMenus[item.ID]; ok {
			for _, actionID := range rms.ToActionIDs() {
				if mrs, ok := mMenuResources[actionID]; ok {
					for _, mr := range mrs {
						if mr.Path == "" || mr.Method == "" {
							continue
						} else if _, ok := mcache[mr.Path+mr.Method]; ok {
							continue
						}
						mcache[mr.Path+mr.Method] = struct{}{}
						line := fmt.Sprintf("p,%d,%s,%s", item.ID, mr.Path, mr.Method)
						persist.LoadPolicyLine(line, m)
					}
				}
			}
		}
	}

	return nil
}

// Load user policy (g,user_id,role_id)
func (a *CasbinAdapter) loadUserPolicy(ctx context.Context, m casbinModel.Model) error {
	userResult, err := a.UserRepo.Query(ctx, schema.UserQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(userResult.Data) > 0 {
		userRoleResult, err := a.UserRoleRepo.Query(ctx, schema.UserRoleQueryParam{})
		if err != nil {
			return err
		}

		mUserRoles := userRoleResult.Data.ToUserIDMap()
		for _, uitem := range userResult.Data {
			if urs, ok := mUserRoles[uitem.ID]; ok {
				for _, ur := range urs {
					line := fmt.Sprintf("g,%d,%d", ur.UserID, ur.RoleID)
					persist.LoadPolicyLine(line, m)
				}
			}
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
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
