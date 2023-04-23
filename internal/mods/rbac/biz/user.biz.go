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

// User management for RBAC
type User struct {
	Trans   *utils.Trans
	UserDAL *dal.User
}

// Query users from the data access object based on the provided parameters and options.
func (a *User) Query(ctx context.Context, params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	params.Pagination = true

	result, err := a.UserDAL.Query(ctx, params, schema.UserQueryOptions{
		QueryOptions: utils.QueryOptions{
			OrderFields: []utils.OrderByParam{
				{Field: "created_at", Direction: utils.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get the specified user from the data access object.
func (a *User) Get(ctx context.Context, id string) (*schema.User, error) {
	user, err := a.UserDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound("", "User not found")
	}
	return user, nil
}

// Create a new user in the data access object.
func (a *User) Create(ctx context.Context, formItem *schema.UserForm) (*schema.User, error) {
	user := &schema.User{
		ID:        idx.NewXID(),
		CreatedAt: time.Now(),
	}
	formItem.FillTo(user)

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.Create(ctx, user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update the specified user in the data access object.
func (a *User) Update(ctx context.Context, id string, formItem *schema.UserForm) error {
	user, err := a.UserDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "User not found")
	}
	formItem.FillTo(user)
	user.UpdatedAt = time.Now()

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.Update(ctx, user); err != nil {
			return err
		}
		return nil
	})
}

// Delete the specified user from the data access object.
func (a *User) Delete(ctx context.Context, id string) error {
	exists, err := a.UserDAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "User not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	})
}
