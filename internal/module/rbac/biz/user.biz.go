package biz

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/util/hash"
	"github.com/LyricTian/gin-admin/v9/pkg/util/xid"
	"github.com/LyricTian/gin-admin/v9/pkg/x/cachex"
	"go.uber.org/zap"
)

type UserBiz struct {
	TransRepo    utilx.TransRepo
	UserRepo     dao.UserRepo
	UserRoleRepo dao.UserRoleRepo
	RoleRepo     dao.RoleRepo
	Cache        cachex.Cacher
}

func (a *UserBiz) Query(ctx context.Context, params typed.UserQueryParam) (*typed.UserQueryResult, error) {
	params.Pagination = true
	result, err := a.UserRepo.Query(ctx, params, typed.UserQueryOptions{
		QueryOptions: utilx.QueryOptions{
			OrderFields: []utilx.OrderByParam{
				{Field: "created_at", Direction: utilx.DESC},
			},
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, nil
	}

	// Fill user roles (include role name)
	userRoleResult, err := a.UserRoleRepo.Query(ctx, typed.UserRoleQueryParam{
		UserIDList: result.Data.ToIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(userRoleResult.Data) > 0 {
		roleResult, err := a.RoleRepo.Query(ctx, typed.RoleQueryParam{
			IDList: userRoleResult.Data.ToRoleIDs(),
			Status: typed.RoleStatusEnabled,
		}, typed.RoleQueryOptions{
			QueryOptions: utilx.QueryOptions{
				SelectFields: []string{"id", "name"},
			},
		})
		if err != nil {
			return nil, err
		}

		userRoleResult.Data.FillRole(roleResult.Data.ToMap())
		result.Data.FillUserRoles(userRoleResult.Data.ToUserIDMap())
	}

	return result, nil
}

func (a *UserBiz) Get(ctx context.Context, id string) (*typed.User, error) {
	user, err := a.UserRepo.Get(ctx, id, typed.UserQueryOptions{
		QueryOptions: utilx.QueryOptions{
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "User not found")
	}

	userRoleResult, err := a.UserRoleRepo.Query(ctx, typed.UserRoleQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	user.UserRoles = userRoleResult.Data

	return user, nil
}

func (a *UserBiz) Create(ctx context.Context, createItem typed.UserCreate) (*typed.User, error) {
	if createItem.Password == "" {
		return nil, errors.BadRequest(errors.ErrBadRequestID, "Password is required")
	}

	exists, err := a.UserRepo.ExistsUsername(ctx, createItem.Username)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest(errors.ErrBadRequestID, "Username already exists")
	}

	hashPass, err := hash.GeneratePassword(createItem.Password)
	if err != nil {
		return nil, errors.BadRequest(errors.ErrBadRequestID, "Failed to generate hash password: %s", err.Error())
	}

	user := &typed.User{
		ID:        xid.NewID(),
		Username:  createItem.Username,
		Password:  hashPass,
		Name:      createItem.Name,
		Phone:     &createItem.Phone,
		Email:     &createItem.Email,
		Remark:    &createItem.Remark,
		Status:    typed.UserStatusActivated,
		CreatedBy: contextx.FromUserID(ctx),
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Create(ctx, user); err != nil {
			return err
		}

		for _, roleID := range createItem.RoleIDs {
			err := a.UserRoleRepo.Create(ctx, &typed.UserRole{
				ID:     xid.NewID(),
				UserID: user.ID,
				RoleID: roleID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserBiz) Update(ctx context.Context, id string, createItem typed.UserCreate) error {
	oldUser, err := a.UserRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldUser == nil {
		return errors.NotFound(errors.ErrNotFoundID, "User not found")
	} else if oldUser.Username != createItem.Username {
		exists, err := a.UserRepo.ExistsUsername(ctx, createItem.Username)
		if err != nil {
			return err
		} else if exists {
			return errors.BadRequest(errors.ErrBadRequestID, "Username already exists")
		}
	}

	// If password is not empty, generate new hash password
	if createItem.Password != "" {
		hashPass, err := hash.GeneratePassword(createItem.Password)
		if err != nil {
			return errors.BadRequest(errors.ErrBadRequestID, "Failed to generate hash password: %s", err.Error())
		}
		oldUser.Password = hashPass
	}

	oldUser.Username = createItem.Username
	oldUser.Name = createItem.Name
	oldUser.Phone = &createItem.Phone
	oldUser.Email = &createItem.Email
	oldUser.Remark = &createItem.Remark
	oldUser.UpdatedBy = contextx.FromUserID(ctx)

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Update(ctx, oldUser); err != nil {
			return err
		}

		if err := a.UserRoleRepo.DeleteByUserID(ctx, id); err != nil {
			return err
		}

		for _, roleID := range createItem.RoleIDs {
			err := a.UserRoleRepo.Create(ctx, &typed.UserRole{
				ID:     xid.NewID(),
				UserID: oldUser.ID,
				RoleID: roleID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Delete user cache by user id
	err = a.Cache.Delete(ctx, utilx.CacheNSForUserRole, id)
	if err != nil {
		logger.Context(ctx).Error("Failed to delete user cache", zap.String("user_id", id), zap.Error(err))
	}

	return nil
}

func (a *UserBiz) Delete(ctx context.Context, id string) error {
	exists, err := a.UserRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound(errors.ErrNotFoundID, "User not found")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Delete(ctx, id); err != nil {
			return err
		}

		if err := a.UserRoleRepo.DeleteByUserID(ctx, id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Delete user cache by user id
	err = a.Cache.Delete(ctx, utilx.CacheNSForUserRole, id)
	if err != nil {
		logger.Context(ctx).Error("Failed to delete user cache", zap.String("user_id", id), zap.Error(err))
	}

	return nil
}

func (a *UserBiz) UpdateStatus(ctx context.Context, id string, status typed.UserStatus) error {
	exists, err := a.UserRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound(errors.ErrNotFoundID, "User not found")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.UserRepo.UpdateStatus(ctx, id, status)
	})
	if err != nil {
		return err
	}

	// Delete user cache by user id
	err = a.Cache.Delete(ctx, utilx.CacheNSForUserRole, id)
	if err != nil {
		logger.Context(ctx).Error("Failed to delete user cache", zap.String("user_id", id), zap.Error(err))
	}

	return nil
}

func (a *UserBiz) GetRoleIDs(ctx context.Context, id string) ([]string, error) {
	user, err := a.UserRepo.Get(ctx, id, typed.UserQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"status"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "User not found")
	} else if user.Status != typed.UserStatusActivated {
		return nil, utilx.ErrUserFreezed
	}

	roleIDs, err := a.UserRoleRepo.GetRoleIDsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}
