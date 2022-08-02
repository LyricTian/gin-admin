package biz

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/jwtauth"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/util/hash"
	"github.com/LyricTian/gin-admin/v9/pkg/x/cachex"
	"go.uber.org/zap"
)

type LoginBiz struct {
	Auth           jwtauth.Auther
	Cache          cachex.Cacher
	UserRepo       *dao.UserRepo
	UserRoleRepo   *dao.UserRoleRepo
	RoleRepo       *dao.RoleRepo
	MenuRepo       *dao.MenuRepo
	MenuActionRepo *dao.MenuActionRepo
}

func (a *LoginBiz) GetCaptchaID(ctx context.Context) (*typed.Captcha, error) {
	item := &typed.Captcha{
		CaptchaID: captcha.NewLen(config.C.Util.Captcha.Length),
	}
	return item, nil
}

func (a *LoginBiz) WriteCaptchaImage(ctx context.Context, w http.ResponseWriter, captchaID string, reload bool) error {
	if reload && !captcha.Reload(captchaID) {
		return errors.NotFound(errors.ErrNotFoundID, "Captcha id not found")
	}

	err := captcha.WriteImage(w, captchaID, config.C.Util.Captcha.Width, config.C.Util.Captcha.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.NotFound(errors.ErrNotFoundID, "Captcha id not found")
		}
		return err
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *LoginBiz) generateToken(ctx context.Context, userID string) (*typed.LoginToken, error) {
	tokenInfo, err := a.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	logger.Context(ctx).Info("Generate token", zap.String("token", tokenInfo.GetAccessToken()), zap.String("user_id", userID))

	return &typed.LoginToken{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}, nil
}

func (a *LoginBiz) Login(ctx context.Context, params typed.UserLogin) (*typed.LoginToken, error) {
	if !captcha.VerifyString(params.CaptchaID, params.CaptchaCode) {
		return nil, errors.BadRequest(errors.ErrBadRequestID, "Invalid captcha code")
	}

	ctx = logger.NewTag(ctx, logger.TagKeyLogin)
	if params.LoginName == config.C.Dictionary.RootUser.Username {
		if params.Password != hash.MD5String(config.C.Dictionary.RootUser.Password) {
			return nil, utilx.ErrInvalidUsernameOrPassword
		}
		return a.generateToken(ctx, config.C.Dictionary.RootUser.ID)
	}

	userResult, err := a.UserRepo.Query(ctx, typed.UserQueryParam{
		Username: params.LoginName,
		PaginationParam: utilx.PaginationParam{
			PageSize: 1,
		},
	}, typed.UserQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id", "password", "status"},
		},
	})
	if err != nil {
		return nil, err
	} else if len(userResult.Data) == 0 {
		return nil, utilx.ErrInvalidUsernameOrPassword
	} else if userResult.Data[0].Status != typed.UserStatusActivated {
		return nil, utilx.ErrUserFreezed
	}

	userID := userResult.Data[0].ID
	ctx = logger.NewUserID(ctx, userID)

	err = hash.CompareHashAndPassword(userResult.Data[0].Password, params.Password)
	if err != nil {
		logger.Context(ctx).Error("Invalid password", zap.Error(err))
		return nil, utilx.ErrInvalidUsernameOrPassword
	}

	// Set user role in cache for permission check
	roleIDs, err := a.UserRoleRepo.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = a.Cache.Set(ctx, utilx.CacheNSForUserRole, userID, strings.Join(roleIDs, ","), time.Hour*time.Duration(config.C.Dictionary.UserCacheExpire))
	if err != nil {
		logger.Context(ctx).Error("Failed to set user role in cache", zap.Error(err))
	}

	logger.Context(ctx).Info("Login success", zap.String("username", params.LoginName))
	return a.generateToken(ctx, userID)
}

func (a *LoginBiz) Logout(ctx context.Context, token string) error {
	ctx = logger.NewTag(ctx, logger.TagKeyLogout)

	err := a.Auth.DestroyToken(ctx, token)
	if err != nil {
		return err
	}

	// Clear user role in cache
	err = a.Cache.Delete(ctx, utilx.CacheNSForUserRole, contextx.FromUserID(ctx))
	if err != nil {
		logger.Context(ctx).Error("Failed to clear user role in cache", zap.Error(err))
	}

	logger.Context(ctx).Info("Logout system")
	return nil
}

func (a *LoginBiz) RefreshToken(ctx context.Context) (*typed.LoginToken, error) {
	userID := contextx.FromUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, typed.UserQueryOptions{
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

	loginToken, err := a.generateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	logger.Context(ctx).Info("Refresh token", zap.String("token", loginToken.AccessToken))
	return loginToken, nil
}

func (a *LoginBiz) GetCurrentUser(ctx context.Context) (*typed.User, error) {
	if utilx.IsRootUser(ctx) {
		return &typed.User{
			ID:        config.C.Dictionary.RootUser.ID,
			Username:  config.C.Dictionary.RootUser.Username,
			Name:      config.C.Dictionary.RootUser.Name,
			Status:    typed.UserStatusActivated,
			CreatedAt: time.Now(),
		}, nil
	}

	userID := contextx.FromUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, typed.UserQueryOptions{
		QueryOptions: utilx.QueryOptions{
			OmitFields: []string{"password", "created_by", "updated_by"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "User not found")
	} else if user.Status != typed.UserStatusActivated {
		return nil, utilx.ErrUserFreezed
	}

	userRoleResult, err := a.UserRoleRepo.Query(ctx, typed.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

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
	user.UserRoles = userRoleResult.Data

	return user, nil
}

func (a *LoginBiz) QueryPrivilegeMenus(ctx context.Context) (typed.Menus, error) {
	userID := contextx.FromUserID(ctx)

	var menuIDList []string
	if !utilx.IsRootUser(ctx) {
		menuResult, err := a.MenuRepo.Query(ctx, typed.MenuQueryParam{
			UserID: userID,
			Status: typed.MenuStatusEnabled,
		}, typed.MenuQueryOptions{
			QueryOptions: utilx.QueryOptions{
				SelectFields: []string{"id", "parent_path"},
			},
		})
		if err != nil {
			return nil, err
		}
		menuIDList = menuResult.Data.SplitParentIDs()
		if len(menuIDList) == 0 {
			return nil, nil
		}
	}

	menuResult, err := a.MenuRepo.Query(ctx, typed.MenuQueryParam{
		IDList: menuIDList,
		Status: typed.MenuStatusEnabled,
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, nil
	}

	menuActionQueryParams := typed.MenuActionQueryParam{}
	if !utilx.IsRootUser(ctx) {
		menuActionQueryParams.UserID = userID
	}
	menuActionResult, err := a.MenuActionRepo.Query(ctx, menuActionQueryParams)
	if err != nil {
		return nil, err
	}
	menuResult.Data.FillActions(menuActionResult.Data.ToMenuIDMap())

	return menuResult.Data.ToTree(), nil
}

func (a *LoginBiz) UpdatePassword(ctx context.Context, params typed.LoginPasswordUpdate) error {
	if utilx.IsRootUser(ctx) {
		return errors.Forbidden(errors.ErrForbiddenID, "Root user can not update password")
	}

	userID := contextx.FromUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, typed.UserQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"id", "password", "status"},
		},
	})
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound(errors.ErrNotFoundID, "User not found")
	} else if user.Status != typed.UserStatusActivated {
		return utilx.ErrUserFreezed
	}

	err = hash.CompareHashAndPassword(user.Password, params.OldPassword)
	if err != nil {
		logger.Context(ctx).Error("Failed to compare old password", zap.Error(err))
		return errors.BadRequest(errors.ErrBadRequestID, "Invalid old password")
	}

	newPassword, err := hash.GeneratePassword(params.NewPassword)
	if err != nil {
		return err
	}

	return a.UserRepo.UpdatePassword(ctx, userID, newPassword)
}
