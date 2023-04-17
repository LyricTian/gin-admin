package biz

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/dal"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/idx"
)

// Role management for RBAC
type Role struct {
	Trans   *utils.Trans
	RoleDAL *dal.Role
}

// Query roles from the data access object based on the provided parameters and options.
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam) (*schema.RoleQueryResult, error) {
	params.Pagination = true

	result, err := a.RoleDAL.Query(ctx, params, schema.RoleQueryOptions{
		QueryOptions: utils.QueryOptions{
			OrderFields: []utils.OrderByParam{
				{Field: "sequence", Direction: utils.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get the specified role from the data access object.
func (a *Role) Get(ctx context.Context, id string) (*schema.Role, error) {
	role, err := a.RoleDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if role == nil {
		return nil, errors.NotFound("", "Role not found")
	}
	return role, nil
}

// Create a new role in the data access object.
func (a *Role) Create(ctx context.Context, formItem *schema.RoleForm) (*schema.Role, error) {
	role := &schema.Role{
		ID:        idx.NewXID(),
		CreatedAt: time.Now(),
	}
	formItem.FillTo(role)

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleDAL.Create(ctx, role); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return role, nil
}

// Update the specified role in the data access object.
func (a *Role) Update(ctx context.Context, id string, formItem *schema.RoleForm) error {
	role, err := a.RoleDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if role == nil {
		return errors.NotFound("", "Role not found")
	}
	formItem.FillTo(role)
	role.UpdatedAt = time.Now()

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleDAL.Update(ctx, role); err != nil {
			return err
		}
		return nil
	})
}

// Delete the specified role from the data access object.
func (a *Role) Delete(ctx context.Context, id string) error {
	exists, err := a.RoleDAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "Role not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleDAL.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	})
}
